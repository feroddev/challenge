package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"

	cfg "github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/handlers"
	"github/feroddev/challengeV3/internal/middleware"
	"github/feroddev/challengeV3/internal/pkg/cache"
	"github/feroddev/challengeV3/internal/pkg/database"
	"github/feroddev/challengeV3/internal/pkg/logger"
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"github/feroddev/challengeV3/internal/pkg/recognition"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
)

func main() {
	loadEnv()

	config := cfg.NewConfig()

	// Inicializa o logger
	isDebug := config.Server.Environment == "development"
	zapLogger, err := logger.NewLogger(isDebug)
	if err != nil {
		log.Fatalf("Erro ao inicializar logger: %v", err)
	}
	defer zapLogger.Sync()

	// Conexão com o banco de dados
	db, err := database.NewDatabaseConnection()
	if err != nil {
		zapLogger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}

	err = database.MigrateDatabase(db)
	if err != nil {
		zapLogger.Fatal("Erro ao migrar banco de dados", zap.Error(err))
	}

	// Inicializa o serviço Redis
	redisCache := cache.NewRedisCache(
		config.Redis.Addr,
		config.Redis.Password,
		config.Redis.DB,
		zapLogger,
	)

	// Inicializa o serviço de reconhecimento
	// Cria o cliente Rekognition diretamente
	rekognitionClient := rekognition.New(rekognition.Options{
		Region: config.AWS.Region,
		Credentials: recognition.NewStaticCredentialsProvider(
			config.AWS.AccessKeyID,
			config.AWS.SecretAccessKey,
			"",
		),
	})
	rekognitionService := recognition.NewRekognitionService(rekognitionClient, zapLogger)

	// Inicializa o produtor NATS
	producerConfig := messaging.ProducerConfig{
		URL:               config.NATS.URL,
		RetryAttempts:     config.NATS.RetryAttempts,
		RetryDelaySeconds: config.NATS.RetryDelaySeconds,
		DeadLetterTopic:   config.NATS.DeadLetterTopic,
	}

	producer, err := messaging.NewProducer(producerConfig, zapLogger)
	if err != nil {
		zapLogger.Fatal("Erro ao inicializar produtor NATS", zap.Error(err))
	}
	defer producer.Close()

	router := setupRouter(config, db, zapLogger, redisCache, rekognitionService, producer)

	zapLogger.Info("Servidor iniciado", zap.String("porta", config.Server.Port))
	log.Fatal(router.Run(":" + config.Server.Port))
}

func setupRouter(config *cfg.Config, db *gorm.DB, zapLogger *zap.Logger, redisCache cache.RedisCache, rekognitionService recognition.RekognitionService, producer *messaging.Producer) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(zapLogger))

	telemetryRepository := repositories.NewPostgresTelemetryRepository(db)
	
	telemetryServiceConfig := services.TelemetryServiceConfig{
		GyroscopeTopic: config.NATS.GyroscopeTopic,
		GPSTopic:       config.NATS.GPSTopic,
		PhotoTopic:     config.NATS.PhotoTopic,
	}
	
	telemetryService := services.NewTelemetryService(
		telemetryRepository,
		producer,
		telemetryServiceConfig,
		zapLogger,
	)
	
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService, zapLogger)

	telemetryGroup := router.Group("/telemetry")
	{
		telemetryGroup.POST("/gyroscope", telemetryHandler.HandleGyroscopeData)
		telemetryGroup.POST("/gps", telemetryHandler.HandleGPSData)
		telemetryGroup.POST("/photo", telemetryHandler.HandlePhotoData)
		
		telemetryGroup.GET("/gyroscope", telemetryHandler.GetGyroscopeData)
		telemetryGroup.GET("/gps", telemetryHandler.GetGPSData)
		telemetryGroup.GET("/photo", telemetryHandler.GetPhotoData)
	}

	return router
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
}
