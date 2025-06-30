package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github/feroddev/challengeV3/internal/configs"
	"github/feroddev/challengeV3/internal/handlers"
	"github/feroddev/challengeV3/internal/middleware"
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"github/feroddev/challengeV3/internal/pkg/metrics"
	"github/feroddev/challengeV3/internal/pkg/tracing"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
)

func main() {
	config := configs.NewConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Inicializa OpenTelemetry
	tracingConfig := tracing.Config{
		ServiceName:    "telemetry-api",
		ServiceVersion: "1.0.0",
		Environment:    config.Environment,
		OTLPEndpoint:   config.Tracing.Endpoint,
	}

	shutdownTracing, err := tracing.InitTracer(tracingConfig, logger)
	if err != nil {
		logger.Error("Falha ao inicializar OpenTelemetry", zap.Error(err))
	} else {
		defer shutdownTracing()
	}

	// Configura o Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Configura middleware de auditoria
	router.Use(middleware.AuditLog(logger))
	
	// Configura middleware de métricas e tracing
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.OpenTelemetryMiddleware("telemetry-api"))
	
	// Configura rate limiter
	rateLimiter := middleware.NewRateLimiter(
		rate.Limit(float64(config.RateLimit.Requests)/float64(config.RateLimit.Period)),
		config.RateLimit.Requests,
		logger,
	)
	
	// Inicia tarefa de limpeza do rate limiter a cada 1 hora
	rateLimiter.CleanupTask(1 * time.Hour)

	// Configura o repositório
	repo, err := repositories.NewTelemetryRepository(config.Database, logger)
	if err != nil {
		logger.Fatal("Falha ao conectar ao banco de dados", zap.Error(err))
	}

	// Configura o produtor NATS com suporte a criptografia
	config.NATS.EncryptionKey = config.Crypto.EncryptionKey
	producer, err := messaging.NewProducer(config.NATS, logger)
	if err != nil {
		logger.Fatal("Falha ao conectar ao NATS", zap.Error(err))
	}
	defer producer.Close()

	// Cria o adaptador para o produtor
	producerAdapter := messaging.NewProducerAdapter(producer)

	// Configura o serviço de telemetria
	telemetryServiceConfig := services.TelemetryServiceConfig{
		GyroscopeTopic: config.Topics.Gyroscope,
		GPSTopic:       config.Topics.GPS,
		PhotoTopic:     config.Topics.Photo,
	}
	telemetryService := services.NewTelemetryService(repo, producerAdapter, telemetryServiceConfig, logger)

	// Configura o serviço de autenticação
	authServiceConfig := services.AuthServiceConfig{
		JWTSecret:     config.Auth.JWTSecret,
		TokenDuration: time.Duration(config.Auth.TokenDuration) * time.Hour,
	}
	db := repositories.GetDB()
	authService := services.NewAuthService(authServiceConfig, db, logger)

	// Configura os handlers
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)

	// Configura as rotas públicas (sem autenticação)
	public := router.Group("/api")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
	
	// Configura as rotas protegidas (com autenticação)
	v1 := router.Group("/api/v1")
	{
		// Aplica middleware de autenticação JWT
		v1.Use(middleware.JWTAuth(authService, logger))
		
		// Aplica rate limiting por IP e Device ID
		v1.Use(rateLimiter.IPRateLimit())
		
		telemetry := v1.Group("/telemetry")
		{
			// Rotas POST com rate limiting por dispositivo
			postGroup := telemetry.Group("/")
			postGroup.Use(rateLimiter.DeviceRateLimit())
			postGroup.POST("/gyroscope", telemetryHandler.HandleGyroscopeData)
			postGroup.POST("/gps", telemetryHandler.HandleGPSData)
			postGroup.POST("/photo", telemetryHandler.HandlePhotoData)
			
			// Rotas GET apenas para usuários com role admin
			getGroup := telemetry.Group("/")
			getGroup.Use(middleware.RoleRequired("admin"))
			getGroup.GET("/gyroscope", telemetryHandler.GetGyroscopeData)
			getGroup.GET("/gps", telemetryHandler.GetGPSData)
			getGroup.GET("/photo", telemetryHandler.GetPhotoData)
		}
	}

	// Configura rota de métricas do Prometheus (protegida por autenticação)
	metricsGroup := router.Group("/metrics")
	metricsGroup.Use(middleware.JWTAuth(authService, logger))
	metricsGroup.Use(middleware.RoleRequired("admin"))
	metricsGroup.GET("", gin.WrapH(metrics.MetricsHandler()))

	// Configura Swagger
	handlers.SetupSwaggerRoutes(router)

	// Inicia o servidor com graceful shutdown
	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Servidor iniciado na porta " + config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Falha ao iniciar servidor", zap.Error(err))
		}
	}()

	// Configura graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Desligando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Erro ao desligar servidor", zap.Error(err))
	}

	logger.Info("Servidor encerrado com sucesso")
}
