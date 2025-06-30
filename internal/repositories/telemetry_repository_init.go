package repositories

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github/feroddev/challengeV3/internal/configs"
	"github/feroddev/challengeV3/internal/core"
)

func NewTelemetryRepository(config configs.DatabaseConfig, logger *zap.Logger) (TelemetryRepository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Falha ao conectar ao banco de dados", zap.Error(err))
		return nil, err
	}

	// Migração automática das tabelas
	err = db.AutoMigrate(&core.Gyroscope{}, &core.GPS{}, &core.Photo{})
	if err != nil {
		logger.Error("Falha ao migrar tabelas", zap.Error(err))
		return nil, err
	}

	logger.Info("Conexão com banco de dados estabelecida com sucesso")
	return NewPostgresTelemetryRepository(db), nil
}
