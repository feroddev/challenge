package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
)

type RekognitionClient interface {
	CompareFaces(sourceImage []byte, targetImage []byte) (float32, error)
}

type RedisCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

type telemetryServiceWithRecognition struct {
	*telemetryService
	rekognition RekognitionClient
	cache       RedisCache
}

func NewTelemetryServiceWithRecognition(
	repository core.TelemetryRepository,
	producer core.Producer,
	rekognition RekognitionClient,
	cache RedisCache,
	config TelemetryServiceConfig,
	logger *zap.Logger,
) *telemetryServiceWithRecognition {
	baseService := &telemetryService{
		repository: repository,
		producer:   producer,
		config:     config,
		logger:     logger,
	}

	return &telemetryServiceWithRecognition{
		telemetryService: baseService,
		rekognition:      rekognition,
		cache:            cache,
	}
}

func (s *telemetryServiceWithRecognition) SavePhotoData(data core.Photo) error {
	if err := s.repository.SavePhotoData(data); err != nil {
		s.logger.Error("Erro ao salvar dados da foto", zap.Error(err))
		return err
	}

	if err := s.producer.Publish(s.config.PhotoTopic, data); err != nil {
		s.logger.Error("Erro ao publicar dados da foto", zap.Error(err))
	}

	go func() {
		if err := s.ProcessPhotoRecognition(data); err != nil {
			s.logger.Error("Erro ao processar reconhecimento facial", zap.Error(err))
		}
	}()

	return nil
}

func (s *telemetryServiceWithRecognition) ProcessPhotoRecognition(photo core.Photo) error {
	cacheKey := fmt.Sprintf("recognition:%d", photo.ID)
	
	cachedResult, err := s.cache.Get(cacheKey)
	if err == nil {
		result := cachedResult.(*core.RecognitionResult)
		return s.repository.UpdatePhotoRecognition(photo.ID, result.Recognized, result.Similarity)
	}

	if err != core.ErrCacheMiss {
		s.logger.Error("Erro ao acessar cache", zap.Error(err))
	}

	previousPhotos, err := s.repository.GetPhotosByDeviceID(photo.DeviceID)
	if err != nil {
		s.logger.Error("Erro ao buscar fotos anteriores", zap.Error(err))
		return err
	}

	if len(previousPhotos) == 0 {
		s.logger.Info("Nenhuma foto anterior encontrada para comparação")
		return nil
	}

	var maxSimilarity float32
	var recognized bool

	for _, prevPhoto := range previousPhotos {
		if prevPhoto.ID == photo.ID {
			continue
		}

		sourceBytes, err := base64.StdEncoding.DecodeString(prevPhoto.Photo)
		if err != nil {
			s.logger.Error("Erro ao decodificar foto de origem", zap.Error(err))
			continue
		}
		
		targetBytes, err := base64.StdEncoding.DecodeString(photo.Photo)
		if err != nil {
			s.logger.Error("Erro ao decodificar foto alvo", zap.Error(err))
			continue
		}
		
		similarity, err := s.rekognition.CompareFaces(sourceBytes, targetBytes)
		if err != nil {
			s.logger.Error("Erro ao comparar faces", zap.Error(err))
			continue
		}

		if similarity > maxSimilarity {
			maxSimilarity = similarity
		}
	}

	recognized = maxSimilarity >= 0.80

	result := &core.RecognitionResult{
		Recognized: recognized,
		Similarity: maxSimilarity,
	}

	s.cache.Set(cacheKey, result, 24*time.Hour)

	return s.repository.UpdatePhotoRecognition(photo.ID, recognized, maxSimilarity)
}
