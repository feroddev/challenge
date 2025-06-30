package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ipLimiters    map[string]*rate.Limiter
	deviceLimiters map[string]*rate.Limiter
	mu            sync.Mutex
	rate          rate.Limit
	burst         int
	logger        *zap.Logger
}

func NewRateLimiter(r rate.Limit, b int, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		ipLimiters:     make(map[string]*rate.Limiter),
		deviceLimiters: make(map[string]*rate.Limiter),
		rate:           r,
		burst:          b,
		logger:         logger,
	}
}

func (rl *RateLimiter) getIPLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ipLimiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ipLimiters[ip] = limiter
	}

	return limiter
}

func (rl *RateLimiter) getDeviceLimiter(deviceID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.deviceLimiters[deviceID]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.deviceLimiters[deviceID] = limiter
	}

	return limiter
}

func (rl *RateLimiter) IPRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getIPLimiter(ip)
		
		if !limiter.Allow() {
			rl.logger.Warn("Taxa de requisicoes excedida por IP",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path))
				
			c.JSON(http.StatusTooManyRequests, gin.H{
				"erro": "taxa de requisicoes excedida, tente novamente mais tarde",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

func (rl *RateLimiter) DeviceRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deviceID string
		
		// Tenta obter o device ID do corpo da requisição
		if c.Request.Method == "POST" {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err == nil {
				if id, exists := body["device_id"]; exists {
					if deviceIDStr, ok := id.(string); ok {
						deviceID = deviceIDStr
					}
				}
			}
			
			// Redefine o corpo da requisição para que possa ser lido novamente
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		}
		
		// Se não encontrou no corpo, tenta obter do parâmetro de consulta
		if deviceID == "" {
			deviceID = c.Query("device_id")
		}
		
		// Se ainda não encontrou, usa o IP como fallback
		if deviceID == "" {
			deviceID = c.ClientIP()
		}
		
		limiter := rl.getDeviceLimiter(deviceID)
		
		if !limiter.Allow() {
			rl.logger.Warn("Taxa de requisicoes excedida por dispositivo",
				zap.String("device_id", deviceID),
				zap.String("path", c.Request.URL.Path))
				
			c.JSON(http.StatusTooManyRequests, gin.H{
				"erro": "taxa de requisicoes excedida para este dispositivo, tente novamente mais tarde",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

func (rl *RateLimiter) CleanupTask(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.cleanup()
		}
	}()
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	// Implementação simples: limpa todos os limitadores
	// Em uma implementação mais sofisticada, poderia verificar quais não foram usados recentemente
	rl.ipLimiters = make(map[string]*rate.Limiter)
	rl.deviceLimiters = make(map[string]*rate.Limiter)
	
	rl.logger.Info("Limpeza dos limitadores de taxa realizada")
}
