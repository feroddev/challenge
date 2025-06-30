package repositories

import (
	"go.uber.org/zap"

	"github/feroddev/challengeV3/internal/configs"
)

func NewTelemetryRepository(config configs.DatabaseConfig, logger *zap.Logger) (TelemetryRepository, error) {
	db, err := InitDB(config, logger)
	if err != nil {
		return nil, err
	}

	return NewPostgresTelemetryRepository(db), nil
}
