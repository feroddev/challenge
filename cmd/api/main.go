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
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.OpenTelemetryMiddleware("telemetry-api"))

	// Configura o repositório
	repo, err := repositories.NewTelemetryRepository(config.Database, logger)
	if err != nil {
		logger.Fatal("Falha ao conectar ao banco de dados", zap.Error(err))
	}

	// Configura o produtor NATS
	producer, err := messaging.NewProducer(config.NATS, logger)
	if err != nil {
		logger.Fatal("Falha ao conectar ao NATS", zap.Error(err))
	}
	defer producer.Close()

	// Configura o serviço de telemetria
	telemetryServiceConfig := services.TelemetryServiceConfig{
		GyroscopeTopic: config.Topics.Gyroscope,
		GPSTopic:       config.Topics.GPS,
		PhotoTopic:     config.Topics.Photo,
	}
	telemetryService := services.NewTelemetryService(repo, producer, telemetryServiceConfig, logger)

	// Configura os handlers
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService, logger)

	// Configura as rotas
	v1 := router.Group("/api/v1")
	{
		telemetry := v1.Group("/telemetry")
		{
			telemetry.POST("/gyroscope", telemetryHandler.HandleGyroscopeData)
			telemetry.POST("/gps", telemetryHandler.HandleGPSData)
			telemetry.POST("/photo", telemetryHandler.HandlePhotoData)
			telemetry.GET("/gyroscope", telemetryHandler.GetGyroscopeData)
			telemetry.GET("/gps", telemetryHandler.GetGPSData)
			telemetry.GET("/photo", telemetryHandler.GetPhotoData)
		}
	}

	// Configura rota de métricas do Prometheus
	router.GET("/metrics", gin.WrapH(metrics.MetricsHandler()))

	// Configura Swagger
	handlers.SetupSwaggerRoutes(router)

	// Configura rota de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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
