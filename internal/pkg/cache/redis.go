package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisCache(addr, password string, db int, logger *zap.Logger) RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisCache{
		client: client,
		logger: logger,
	}
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("Erro ao serializar valor para cache", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	err = r.client.Set(ctx, key, json, expiration).Err()
	if err != nil {
		r.logger.Error("Erro ao definir valor no cache", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	r.logger.Debug("Valor armazenado no cache", 
		zap.String("key", key),
		zap.Duration("expiration", expiration))
	return nil
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		r.logger.Debug("Chave não encontrada no cache", zap.String("key", key))
		return err
	} else if err != nil {
		r.logger.Error("Erro ao obter valor do cache", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		r.logger.Error("Erro ao deserializar valor do cache", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	r.logger.Debug("Valor recuperado do cache", zap.String("key", key))
	return nil
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Erro ao excluir chave do cache", 
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	r.logger.Debug("Chave excluída do cache", zap.String("key", key))
	return nil
}
