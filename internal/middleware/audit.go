package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AuditLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Captura o corpo da requisição
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// Captura a resposta
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		
		// Processa a requisição
		c.Next()
		
		// Extrai informações do usuário autenticado
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		
		// Registra a auditoria
		logger.Info("Auditoria",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Any("user_id", userID),
			zap.Any("username", username),
			zap.Int("status", c.Writer.Status()),
			zap.String("latency", time.Since(start).String()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("request_id", c.GetString("request_id")),
			zap.String("request_body", sanitizeRequestBody(string(requestBody))),
			zap.String("response_body", sanitizeResponseBody(blw.body.String())))
	}
}

func sanitizeRequestBody(body string) string {
	// Implementação simplificada: limita o tamanho do corpo
	// Em uma implementação real, você também removeria dados sensíveis
	maxLen := 1000
	if len(body) > maxLen {
		return body[:maxLen] + "..."
	}
	return body
}

func sanitizeResponseBody(body string) string {
	// Implementação simplificada: limita o tamanho do corpo
	// Em uma implementação real, você também removeria dados sensíveis
	maxLen := 1000
	if len(body) > maxLen {
		return body[:maxLen] + "..."
	}
	return body
}
