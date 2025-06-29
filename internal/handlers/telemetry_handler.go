package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
)

type TelemetryHandler struct {
	telemetryService services.TelemetryService
}

func NewTelemetryHandler(telemetryService services.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryService: telemetryService,
	}
}

func (h *TelemetryHandler) HandleGyroscopeData(c *gin.Context) {
	var gyroscopeData core.Gyroscope
	if err := c.ShouldBindJSON(&gyroscopeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gyroscopeData.Timestamp.IsZero() {
		gyroscopeData.Timestamp = time.Now()
	}

	err := h.telemetryService.SaveGyroscopeData(gyroscopeData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dados do giroscópio recebidos com sucesso"})
}

func (h *TelemetryHandler) HandleGPSData(c *gin.Context) {
	var gpsData core.GPS
	if err := c.ShouldBindJSON(&gpsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gpsData.Timestamp.IsZero() {
		gpsData.Timestamp = time.Now()
	}

	err := h.telemetryService.SaveGPSData(gpsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dados do GPS recebidos com sucesso"})
}

func (h *TelemetryHandler) HandlePhotoData(c *gin.Context) {
	var photoData core.Photo
	if err := c.ShouldBindJSON(&photoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if photoData.Timestamp.IsZero() {
		photoData.Timestamp = time.Now()
	}

	err := h.telemetryService.SavePhotoData(photoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dados da foto recebidos com sucesso"})
}

func (h *TelemetryHandler) GetGyroscopeData(c *gin.Context) {
	data, err := h.telemetryService.GetGyroscopeData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *TelemetryHandler) GetGPSData(c *gin.Context) {
	data, err := h.telemetryService.GetGPSData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *TelemetryHandler) GetPhotoData(c *gin.Context) {
	data, err := h.telemetryService.GetPhotoData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
