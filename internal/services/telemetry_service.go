package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github/feroddev/challengeV3/internal/core"
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
	ProcessPhotoRecognition(ctx context.Context, photo core.Photo) (RecognitionResult, error)
	SetAWSRekognition(aws AWSRekognition)
	SetRedisCache(cache RedisCache)
}

type RecognitionResult struct {
	Recognized bool
	Similarity float32
}

type AWSRekognition interface {
	CompareFaces(sourceImage, targetImage string) (bool, float32, error)
}

type RedisCache interface {
	Get(key string) (string, error)
	Set(key, value string, expiration time.Duration) error
}

type telemetryService struct {
	repository     repositories.TelemetryRepository
	producer       core.Producer
	config         TelemetryServiceConfig
	logger         *zap.Logger
	awsRekognition AWSRekognition
	redisCache     RedisCache
}

type TelemetryServiceConfig struct {
	GyroscopeTopic string
	GPSTopic       string
	PhotoTopic     string
}

func NewTelemetryService(repository repositories.TelemetryRepository, producer core.Producer, config TelemetryServiceConfig, logger *zap.Logger) TelemetryService {
	return &telemetryService{
		repository: repository,
		producer:   producer,
		config:     config,
		logger:     logger,
	}
}

func (s *telemetryService) SetAWSRekognition(aws AWSRekognition) {
	s.awsRekognition = aws
}

func (s *telemetryService) SetRedisCache(cache RedisCache) {
	s.redisCache = cache
}

func (s *telemetryService) ProcessPhotoRecognition(ctx context.Context, photo core.Photo) (RecognitionResult, error) {
	result := RecognitionResult{
		Recognized: false,
		Similarity: 0,
	}

	s.logger.Info("Iniciando processamento de reconhecimento facial", 
		zap.String("deviceID", photo.DeviceID),
		zap.Int64("photoID", int64(photo.ID)))

	photos, err := s.repository.GetPhotosByDeviceID(photo.DeviceID)
	if err != nil {
		s.logger.Error("Erro ao buscar fotos do dispositivo", 
			zap.Error(err), 
			zap.String("deviceID", photo.DeviceID))
		return result, err
	}

	if len(photos) <= 1 {
		s.logger.Info("Sem fotos anteriores para comparação", 
			zap.String("deviceID", photo.DeviceID))
		err = s.repository.UpdatePhotoRecognition(photo.ID, false, 0)
		if err != nil {
			s.logger.Error("Erro ao atualizar reconhecimento da foto", zap.Error(err))
		}
		s.publishPhotoResult(ctx, photo, false, 0)
		return result, nil
	}

	var highestSimilarity float32
	var recognized bool

	for _, prevPhoto := range photos {
		if prevPhoto.ID == photo.ID {
			continue
		}

		cacheKey := fmt.Sprintf("recognition:%d:%d", prevPhoto.ID, photo.ID)
		cachedResult, err := s.redisCache.Get(cacheKey)
		if err == nil && cachedResult != "" {
			s.logger.Info("Resultado encontrado em cache", 
				zap.String("key", cacheKey), 
				zap.String("result", cachedResult))
			
			parts := strings.Split(cachedResult, ":")
			if len(parts) == 2 {
				recognized = parts[0] == "true"
				sim, _ := strconv.ParseFloat(parts[1], 32)
				similarity := float32(sim)
				if similarity > highestSimilarity {
					highestSimilarity = similarity
					recognized = true
				}
			}
			continue
		}

		s.logger.Info("Comparando faces", 
			zap.Int64("photoID1", int64(prevPhoto.ID)), 
			zap.Int64("photoID2", int64(photo.ID)))
		
		match, similarity, err := s.awsRekognition.CompareFaces(prevPhoto.Photo, photo.Photo)
		if err != nil {
			s.logger.Error("Erro ao comparar faces", zap.Error(err))
			continue
		}

		cacheValue := fmt.Sprintf("%t:%f", match, similarity)
		err = s.redisCache.Set(cacheKey, cacheValue, 24*time.Hour)
		if err != nil {
			s.logger.Error("Erro ao salvar resultado em cache", zap.Error(err))
		}

		if match && similarity > highestSimilarity {
			highestSimilarity = similarity
			recognized = true
		}
	}

	s.logger.Info("Resultado do reconhecimento", 
		zap.Bool("recognized", recognized), 
		zap.Float32("similarity", highestSimilarity))

	err = s.repository.UpdatePhotoRecognition(photo.ID, recognized, highestSimilarity)
	if err != nil {
		s.logger.Error("Erro ao atualizar reconhecimento da foto", zap.Error(err))
		return result, err
	}

	result.Recognized = recognized
	result.Similarity = highestSimilarity

	s.publishPhotoResult(ctx, photo, recognized, highestSimilarity)

	return result, nil
}

func (s *telemetryService) publishPhotoResult(ctx context.Context, photo core.Photo, recognized bool, similarity float32) {
	photo.Recognized = recognized
	photo.Similarity = similarity
	
	err := s.producer.PublishWithRetry(ctx, s.config.PhotoTopic, photo, 3, 500*time.Millisecond)
	if err != nil {
		s.logger.Error("Erro ao publicar resultado do reconhecimento", zap.Error(err))
	}
}

func (s *telemetryService) SaveGyroscopeData(data core.Gyroscope) error {
	err := s.repository.SaveGyroscopeData(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.producer.PublishWithRetry(ctx, s.config.GyroscopeTopic, data, 3, 2*time.Second)
	if err != nil {
		s.logger.Error("Erro ao publicar dados do giroscópio no NATS",
			zap.Error(err),
			zap.String("device_id", data.DeviceID))
	}

	return nil
}

func (s *telemetryService) SaveGPSData(data core.GPS) error {
	err := s.repository.SaveGPSData(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.producer.PublishWithRetry(ctx, s.config.GPSTopic, data, 3, 2*time.Second)
	if err != nil {
		s.logger.Error("Erro ao publicar dados do GPS no NATS",
			zap.Error(err),
			zap.String("device_id", data.DeviceID))
	}

	return nil
}

func (s *telemetryService) SavePhotoData(data core.Photo) error {
	err := s.repository.SavePhotoData(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.producer.PublishWithRetry(ctx, s.config.PhotoTopic, data, 3, 2*time.Second)
	if err != nil {
		s.logger.Error("Erro ao publicar dados da foto no NATS",
			zap.Error(err),
			zap.String("device_id", data.DeviceID))
	} else {
		s.logger.Info("Foto enviada para processamento assíncrono",
			zap.String("device_id", data.DeviceID),
			zap.String("topic", s.config.PhotoTopic))
	}

	return nil
}

func (s *telemetryService) GetGyroscopeData() ([]core.Gyroscope, error) {
	return s.repository.GetGyroscopeData()
}

func (s *telemetryService) GetGPSData() ([]core.GPS, error) {
	return s.repository.GetGPSData()
}

func (s *telemetryService) GetPhotoData() ([]core.Photo, error) {
	return s.repository.GetPhotoData()
}
