package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/services"
)

func JWTAuth(authService *services.AuthService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "token nao fornecido"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "formato de token invalido"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Falha na autenticacao", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "token invalido"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		logger.Info("Usuario autenticado",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("role", claims.Role))

		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "usuario nao autenticado"})
			c.Abort()
			return
		}

		roleMatched := false
		for _, role := range roles {
			if role == userRole.(string) {
				roleMatched = true
				break
			}
		}

		if !roleMatched {
			c.JSON(http.StatusForbidden, gin.H{"erro": "permissao negada"})
			c.Abort()
			return
		}

		c.Next()
	}
}
