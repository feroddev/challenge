package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cfg "github/feroddev/challengeV3/config"
	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/pkg/cache"
	"github/feroddev/challengeV3/internal/pkg/logger"
	"github/feroddev/challengeV3/internal/pkg/messaging"
	"github/feroddev/challengeV3/internal/pkg/recognition"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	config := cfg.NewConfig()
	zapLogger, err := logger.NewLogger(config.Server.Environment)
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		zapLogger.Fatal("Erro ao conectar ao Redis", zap.Error(err))
	}

	redisCache := cache.NewRedisCache(redisClient, zapLogger)

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(config.AWS.Region))
	if err != nil {
		zapLogger.Fatal("Erro ao carregar configuração da AWS", zap.Error(err))
	}

	rekognitionClient := rekognition.NewFromConfig(awsCfg)
	recognitionService := recognition.NewRekognitionService(rekognitionClient, zapLogger)

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

	zapLogger.Info("Consumidor de fotos iniciado",
		zap.String("topic", config.NATS.PhotoTopic))

	err = consumer.Subscribe(config.NATS.PhotoTopic, func(data []byte) error {
		var photoData core.Photo
		if err := json.Unmarshal(data, &photoData); err != nil {
			return fmt.Errorf("erro ao deserializar dados da foto: %w", err)
		}

		zapLogger.Info("Processando dados de foto",
			zap.String("device_id", photoData.DeviceID),
			zap.Int("photo_length", len(photoData.Photo)))

		if result := db.Create(&photoData); result.Error != nil {
			return fmt.Errorf("erro ao salvar dados da foto: %w", result.Error)
		}

		zapLogger.Info("Foto salva com sucesso, iniciando reconhecimento facial",
			zap.String("device_id", photoData.DeviceID),
			zap.Uint("id", photoData.ID))

		go processPhotoRecognition(ctx, db, redisCache, recognitionService, photoData, zapLogger)

		return nil
	})

	if err != nil {
		zapLogger.Fatal("Erro ao assinar tópico", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Encerrando consumidor de fotos")
}

func processPhotoRecognition(
	ctx context.Context,
	db *gorm.DB,
	cache cache.RedisCache,
	recognitionService recognition.RekognitionService,
	photo core.Photo,
	logger *zap.Logger,
) {
	cacheKey := fmt.Sprintf("photos:%s", photo.DeviceID)
	
	var previousPhotos []core.Photo
	
	cachedPhotos, err := cache.Get(ctx, cacheKey)
	if err == nil {
		if err := json.Unmarshal([]byte(cachedPhotos), &previousPhotos); err != nil {
			logger.Error("Erro ao deserializar fotos do cache", zap.Error(err))
		}
	}
	
	if len(previousPhotos) == 0 {
		var dbPhotos []core.Photo
		result := db.Where("device_id = ? AND id != ?", photo.DeviceID, photo.ID).
			Order("created_at desc").
			Limit(5).
			Find(&dbPhotos)
		
		if result.Error != nil {
			logger.Error("Erro ao buscar fotos anteriores no banco", zap.Error(result.Error))
			return
		}
		
		previousPhotos = dbPhotos
		
		photosJson, err := json.Marshal(previousPhotos)
		if err == nil {
			cache.Set(ctx, cacheKey, string(photosJson), 30*time.Minute)
		}
	}
	
	if len(previousPhotos) == 0 {
		logger.Info("Nenhuma foto anterior encontrada para comparação",
			zap.String("device_id", photo.DeviceID))
		return
	}
	
	var highestSimilarity float32
	recognized := false
	
	for _, prevPhoto := range previousPhotos {
		similarity, err := recognitionService.CompareFaces(ctx, photo.Data, prevPhoto.Data)
		if err != nil {
			logger.Error("Erro ao comparar faces",
				zap.Error(err),
				zap.Uint("current_photo_id", photo.ID),
				zap.Uint("previous_photo_id", prevPhoto.ID))
			continue
		}
		
		if similarity > highestSimilarity {
			highestSimilarity = similarity
		}
		
		if similarity >= 90.0 {
			recognized = true
			break
		}
	}
	
	logger.Info("Resultado do reconhecimento facial",
		zap.String("device_id", photo.DeviceID),
		zap.Uint("photo_id", photo.ID),
		zap.Bool("recognized", recognized),
		zap.Float32("similarity", highestSimilarity))
	
	db.Model(&photo).Updates(map[string]interface{}{
		"recognized": recognized,
		"similarity": highestSimilarity,
	})
	
	recognitionResult := struct {
		PhotoID    uint    `json:"photo_id"`
		Recognized bool    `json:"recognized"`
		Similarity float32 `json:"similarity"`
	}{
		PhotoID:    photo.ID,
		Recognized: recognized,
		Similarity: highestSimilarity,
	}
	
	resultJson, _ := json.Marshal(recognitionResult)
	cache.Set(ctx, fmt.Sprintf("recognition:%d", photo.ID), string(resultJson), 24*time.Hour)
}
