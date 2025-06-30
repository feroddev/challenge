package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user core.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuario registrado com sucesso"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request core.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	token, err := h.authService.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}
