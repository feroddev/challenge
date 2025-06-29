package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github/feroddev/challengeV3/internal/pkg/metrics"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		
		metrics.RequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Observe(duration)
		
		metrics.RequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Inc()
	}
}
