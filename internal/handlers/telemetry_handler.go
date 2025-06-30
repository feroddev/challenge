package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
)

type TelemetryHandler struct {
	telemetryService services.TelemetryService
	logger          *zap.Logger
}

func NewTelemetryHandler(telemetryService services.TelemetryService, logger *zap.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryService: telemetryService,
		logger:          logger,
	}
}

// HandleGyroscopeData recebe dados do giroscópio e salva no repositório
// @Summary Receber dados do giroscópio
// @Description Recebe e salva dados do giroscópio de dispositivos
// @Tags telemetria
// @Accept json
// @Produce json
// @Param data body core.Gyroscope true "Dados do giroscópio"
// @Success 201 {object} map[string]string "Mensagem de sucesso"
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /telemetry/gyroscope [post]
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

// HandleGPSData recebe dados de GPS e salva no repositório
// @Summary Receber dados de GPS
// @Description Recebe e salva dados de GPS de dispositivos
// @Tags telemetria
// @Accept json
// @Produce json
// @Param data body core.GPS true "Dados de GPS"
// @Success 201 {object} map[string]string "Mensagem de sucesso"
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /telemetry/gps [post]
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

// HandlePhotoData recebe dados de foto e salva no repositório
// @Summary Receber dados de foto
// @Description Recebe e salva dados de foto de dispositivos
// @Tags telemetria
// @Accept json
// @Produce json
// @Param data body core.Photo true "Dados da foto"
// @Success 201 {object} map[string]string "Mensagem de sucesso"
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /telemetry/photo [post]
func (h *TelemetryHandler) HandlePhotoData(c *gin.Context) {
	var photoData core.Photo
	if err := c.ShouldBindJSON(&photoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if photoData.Timestamp.IsZero() {
		photoData.Timestamp = time.Now()
	}

	h.logger.Info("Recebendo dados de foto", 
		zap.String("device_id", photoData.DeviceID),
		zap.Time("timestamp", photoData.Timestamp))

	err := h.telemetryService.SavePhotoData(photoData)
	if err != nil {
		h.logger.Error("Erro ao salvar dados da foto", zap.Error(err))
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
	h.logger.Info("Buscando dados de fotos")

	photos, err := h.telemetryService.GetPhotoData()
	if err != nil {
		h.logger.Error("Erro ao buscar dados de fotos", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Dados de fotos recuperados com sucesso", zap.Int("quantidade", len(photos)))
	c.JSON(http.StatusOK, photos)
}
