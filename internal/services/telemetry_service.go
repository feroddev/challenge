package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/cache"
	"github/feroddev/challengeV3/internal/pkg/recognition"
	"github/feroddev/challengeV3/internal/repositories"

	"go.uber.org/zap"
)

type TelemetryService interface {
	SaveGyroscopeData(data core.Gyroscope) error
	SaveGPSData(data core.GPS) error
	SavePhotoData(data core.Photo) error
	GetGyroscopeData() ([]core.Gyroscope, error)
	GetGPSData() ([]core.GPS, error)
	GetPhotoData() ([]core.Photo, error)
	ProcessPhotoRecognition(ctx context.Context, photo core.Photo) (core.Photo, error)
}

type telemetryService struct {
	repository repositories.TelemetryRepository
	rekognition recognition.RekognitionService
	cache      cache.RedisCache
	logger     *zap.Logger
}

func NewTelemetryService(repository repositories.TelemetryRepository, rekognition recognition.RekognitionService, redisCache cache.RedisCache, logger *zap.Logger) TelemetryService {
	return &telemetryService{
		repository:  repository,
		rekognition: rekognition,
		cache:       redisCache,
		logger:      logger,
	}
}

func (s *telemetryService) SaveGyroscopeData(data core.Gyroscope) error {
	return s.repository.SaveGyroscopeData(data)
}

func (s *telemetryService) SaveGPSData(data core.GPS) error {
	return s.repository.SaveGPSData(data)
}

func (s *telemetryService) SavePhotoData(data core.Photo) error {
	// Primeiro salva a foto no banco de dados
	err := s.repository.SavePhotoData(data)
	if err != nil {
		return err
	}

	// Processa o reconhecimento de forma assíncrona
	go func() {
		ctx := context.Background()
		_, err := s.ProcessPhotoRecognition(ctx, data)
		if err != nil {
			s.logger.Error("Erro ao processar reconhecimento de foto", 
				zap.String("device_id", data.DeviceID),
				zap.Error(err))
		}
	}()

	return nil
}

func (s *telemetryService) GetGyroscopeData() ([]core.Gyroscope, error) {
	return s.repository.GetGyroscopeData()
}

func (s *telemetryService) GetGPSData() ([]core.GPS, error) {
	return s.repository.GetGPSData()
}

func (s *telemetryService) GetPhotoData() ([]core.Photo, error) {
	// Tenta buscar do cache primeiro
	ctx := context.Background()
	var photos []core.Photo
	cacheKey := "photos:all"

	err := s.cache.Get(ctx, cacheKey, &photos)
	if err == nil {
		s.logger.Debug("Dados de fotos recuperados do cache")
		return photos, nil
	}

	// Se não estiver no cache, busca do banco de dados
	photos, err = s.repository.GetPhotoData()
	if err != nil {
		return nil, err
	}

	// Armazena no cache por 5 minutos
	s.cache.Set(ctx, cacheKey, photos, 5*time.Minute)

	return photos, nil
}
