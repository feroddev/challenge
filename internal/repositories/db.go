package repositories

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github/feroddev/challengeV3/internal/configs"
	"github/feroddev/challengeV3/internal/core"
)

var dbInstance *gorm.DB

func InitDB(config configs.DatabaseConfig, logger *zap.Logger) (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

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

	err = db.AutoMigrate(&core.Gyroscope{}, &core.GPS{}, &core.Photo{}, &core.User{})
	if err != nil {
		logger.Error("Falha ao migrar tabelas", zap.Error(err))
		return nil, err
	}

	dbInstance = db
	logger.Info("Conexao com banco de dados estabelecida com sucesso")
	return db, nil
}

func GetDB() *gorm.DB {
	return dbInstance
}
