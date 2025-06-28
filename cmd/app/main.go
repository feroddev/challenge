package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/handlers"
	"github/feroddev/challengeV3/internal/repositories"
	"github/feroddev/challengeV3/internal/services"
)

func main() {
	cfg := config.NewConfig()
	router := setupRouter(cfg)

	fmt.Printf("Servidor iniciado na porta %s\n", cfg.Server.Port)
	log.Fatal(router.Run(":" + cfg.Server.Port))
}

func setupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	telemetryRepository := repositories.NewInMemoryTelemetryRepository()
	telemetryService := services.NewTelemetryService(telemetryRepository)
	telemetryHandler := handlers.NewTelemetryHandler(telemetryService)

	telemetryGroup := router.Group("/telemetry")
	{
		telemetryGroup.POST("/gyroscope", telemetryHandler.HandleGyroscopeData)
		telemetryGroup.POST("/gps", telemetryHandler.HandleGPSData)
		telemetryGroup.POST("/photo", telemetryHandler.HandlePhotoData)
	}

	return router
}
