package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
)

type TelemetryHandler struct {
	telemetryService services.TelemetryService
	validator        *validator.Validate
}

func NewTelemetryHandler(telemetryService services.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryService: telemetryService,
		validator:        validator.New(),
	}
}

func (h *TelemetryHandler) HandleGyroscopeData(c *gin.Context) {
	var gyroscopeData core.Gyroscope
	if err := c.ShouldBindJSON(&gyroscopeData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Erro ao processar dados do giroscópio",
		})
		return
	}

	if err := h.validator.Struct(gyroscopeData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Dados do giroscópio inválidos",
		})
		return
	}

	if err := h.telemetryService.SaveGyroscopeData(gyroscopeData); err != nil {
		c.JSON(http.StatusInternalServerError, core.Response{
			Status:  http.StatusInternalServerError,
			Message: "Erro ao salvar dados do giroscópio",
		})
		return
	}

	c.JSON(http.StatusOK, core.Response{
		Status:  http.StatusOK,
		Message: "Dados do giroscópio recebidos com sucesso",
	})
}

func (h *TelemetryHandler) HandleGPSData(c *gin.Context) {
	var gpsData core.GPS
	if err := c.ShouldBindJSON(&gpsData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Erro ao processar dados do GPS",
		})
		return
	}

	if err := h.validator.Struct(gpsData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Dados do GPS inválidos",
		})
		return
	}

	if err := h.telemetryService.SaveGPSData(gpsData); err != nil {
		c.JSON(http.StatusInternalServerError, core.Response{
			Status:  http.StatusInternalServerError,
			Message: "Erro ao salvar dados do GPS",
		})
		return
	}

	c.JSON(http.StatusOK, core.Response{
		Status:  http.StatusOK,
		Message: "Dados do GPS recebidos com sucesso",
	})
}

func (h *TelemetryHandler) HandlePhotoData(c *gin.Context) {
	var photoData core.Photo
	if err := c.ShouldBindJSON(&photoData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Erro ao processar dados da foto",
		})
		return
	}

	if err := h.validator.Struct(photoData); err != nil {
		c.JSON(http.StatusBadRequest, core.Response{
			Status:  http.StatusBadRequest,
			Message: "Dados da foto inválidos",
		})
		return
	}

	if err := h.telemetryService.SavePhotoData(photoData); err != nil {
		c.JSON(http.StatusInternalServerError, core.Response{
			Status:  http.StatusInternalServerError,
			Message: "Erro ao salvar dados da foto",
		})
		return
	}

	c.JSON(http.StatusOK, core.Response{
		Status:  http.StatusOK,
		Message: "Dados da foto recebidos com sucesso",
	})
}
