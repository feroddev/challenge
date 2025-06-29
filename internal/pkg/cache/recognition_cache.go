package cache

import (
	"context"
	"time"

	"github/feroddev/challengeV3/internal/core"
	"go.uber.org/zap"
)

type RecognitionCache struct {
	redisCache RedisCache
	logger     *zap.Logger
}

func NewRecognitionCache(redisCache RedisCache, logger *zap.Logger) *RecognitionCache {
	return &RecognitionCache{
		redisCache: redisCache,
		logger:     logger,
	}
}

func (c *RecognitionCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	var result core.RecognitionResult
	
	err := c.redisCache.Get(ctx, key, &result)
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, core.ErrCacheMiss
		}
		return nil, err
	}
	
	return &result, nil
}

func (c *RecognitionCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return c.redisCache.Set(ctx, key, value, expiration)
}
