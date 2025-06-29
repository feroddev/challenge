package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Processa a requisição
		c.Next()

		// Após o processamento
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("Requisição HTTP",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Int("size", bodySize),
			zap.Duration("latency", latency),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
		)
	}
}
