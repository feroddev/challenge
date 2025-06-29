package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cfg "github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/handlers"
	"github/feroddev/challengeV3/internal/middleware"
	"github/feroddev/challengeV3/internal/pkg/cache"
	"github/feroddev/challengeV3/internal/pkg/database"
	"github/feroddev/challengeV3/internal/pkg/logger"
	"github/feroddev/challengeV3/internal/pkg/recognition"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
)

func main() {
	loadEnv()

	cfg := config.NewConfig()

	// Inicializa o logger
	zapLogger, err := logger.NewLogger(cfg.Server.Environment == "development")
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
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
		0,
		zapLogger,
	)

	// Inicializa o serviço de reconhecimento
	rekognitionService, err := recognition.NewRekognitionService(zapLogger)
	if err != nil {
		zapLogger.Fatal("Erro ao inicializar serviço de reconhecimento", zap.Error(err))
	}

	router := setupRouter(cfg, db, zapLogger, redisCache, rekognitionService)

	zapLogger.Info("Servidor iniciado", zap.String("porta", cfg.Server.Port))
	log.Fatal(router.Run(":" + cfg.Server.Port))
}

func setupRouter(cfg *config.Config, db *gorm.DB, zapLogger *zap.Logger, redisCache cache.RedisCache, rekognitionService recognition.RekognitionService) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware(zapLogger))

	telemetryRepository := repositories.NewPostgresTelemetryRepository(db)
	telemetryService := services.NewTelemetryService(telemetryRepository, rekognitionService, redisCache, zapLogger)
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
