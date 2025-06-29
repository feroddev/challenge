package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cfg "github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/logger"
	"github/feroddev/challengeV3/internal/pkg/messaging"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	config := cfg.NewConfig()
	isDevelopment := config.Server.Environment == "development"
	zapLogger, err := logger.NewLogger(isDevelopment)
	if err != nil {
		fmt.Printf("Erro ao criar logger: %v\n", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name, config.DB.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zapLogger.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}

	consumerConfig := messaging.ConsumerConfig{
		URL:               config.NATS.URL,
		RetryAttempts:     config.NATS.RetryAttempts,
		RetryDelaySeconds: config.NATS.RetryDelaySeconds,
		DeadLetterTopic:   config.NATS.DeadLetterTopic,
	}

	consumer, err := messaging.NewConsumer(consumerConfig, zapLogger)
	if err != nil {
		zapLogger.Fatal("Erro ao criar consumidor NATS", zap.Error(err))
	}
	defer consumer.Close()

	zapLogger.Info("Consumidor de giroscópio iniciado",
		zap.String("topic", config.NATS.GyroscopeTopic))

	err = consumer.Subscribe(config.NATS.GyroscopeTopic, func(data []byte) error {
		var gyroscopeData core.Gyroscope
		if err := json.Unmarshal(data, &gyroscopeData); err != nil {
			return fmt.Errorf("erro ao deserializar dados do giroscópio: %w", err)
		}

		zapLogger.Info("Processando dados de giroscópio",
			zap.String("device_id", gyroscopeData.DeviceID),
			zap.Float64("x", gyroscopeData.X),
			zap.Float64("y", gyroscopeData.Y),
			zap.Float64("z", gyroscopeData.Z))

		if result := db.Create(&gyroscopeData); result.Error != nil {
			return fmt.Errorf("erro ao salvar dados do giroscópio: %w", result.Error)
		}

		zapLogger.Info("Dados de giroscópio salvos com sucesso",
			zap.String("device_id", gyroscopeData.DeviceID))

		return nil
	})

	if err != nil {
		zapLogger.Fatal("Erro ao assinar tópico", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Encerrando consumidor de giroscópio")
}
