package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"github/feroddev/challengeV3/internal/repositories"
)

type TelemetryService interface {
	SaveGyroscopeData(data core.Gyroscope) error
	SaveGPSData(data core.GPS) error
	SavePhotoData(data core.Photo) error
	GetGyroscopeData() ([]core.Gyroscope, error)
	GetGPSData() ([]core.GPS, error)
	GetPhotoData() ([]core.Photo, error)
}

type telemetryService struct {
	repository repositories.TelemetryRepository
	producer   *messaging.Producer
	config     TelemetryServiceConfig
	logger     *zap.Logger
}

type TelemetryServiceConfig struct {
	GyroscopeTopic string
	GPSTopic       string
	PhotoTopic     string
}

func NewTelemetryService(
	repository repositories.TelemetryRepository,
	producer *messaging.Producer,
	config TelemetryServiceConfig,
	logger *zap.Logger,
) TelemetryService {
	return &telemetryService{
		repository: repository,
		producer:   producer,
		config:     config,
		logger:     logger,
	}
}

func (s *telemetryService) SaveGyroscopeData(data core.Gyroscope) error {
	err := s.repository.SaveGyroscopeData(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.producer.PublishWithRetry(ctx, s.config.GyroscopeTopic, data)
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

	err = s.producer.PublishWithRetry(ctx, s.config.GPSTopic, data)
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

	err = s.producer.PublishWithRetry(ctx, s.config.PhotoTopic, data)
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
