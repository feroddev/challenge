package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/handlers"
	"github/feroddev/challengeV3/internal/pkg/database"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
)

func main() {
	loadEnv()

	cfg := config.NewConfig()

	db, err := database.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	err = database.MigrateDatabase(db)
	if err != nil {
		log.Fatalf("Erro ao migrar banco de dados: %v", err)
	}

	router := setupRouter(cfg, db)

	fmt.Printf("Servidor iniciado na porta %s\n", cfg.Server.Port)
	log.Fatal(router.Run(":" + cfg.Server.Port))
}

func setupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	telemetryRepository := repositories.NewPostgresTelemetryRepository(db)
	telemetryService := services.NewTelemetryService(telemetryRepository)
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService)

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
