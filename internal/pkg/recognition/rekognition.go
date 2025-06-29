package recognition

import (
	"context"
	"encoding/base64"
	"errors"
	"github/feroddev/challengeV3/internal/core"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"go.uber.org/zap"
)

type RekognitionService interface {
	CompareFaces(ctx context.Context, sourceImage, targetImage []byte) (bool, float32, error)
	CompareWithPreviousPhotos(ctx context.Context, photo core.Photo) (bool, float32, error)
}

type rekognitionService struct {
	client *rekognition.Client
	logger *zap.Logger
}

func NewRekognitionService(logger *zap.Logger) (RekognitionService, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := rekognition.NewFromConfig(cfg)

	return &rekognitionService{
		client: client,
		logger: logger,
	}, nil
}

func (r *rekognitionService) CompareFaces(ctx context.Context, sourceImage, targetImage []byte) (bool, float32, error) {
	input := &rekognition.CompareFacesInput{
		SourceImage: &types.Image{
			Bytes: sourceImage,
		},
		TargetImage: &types.Image{
			Bytes: targetImage,
		},
		SimilarityThreshold: aws.Float32(70.0),
	}

	output, err := r.client.CompareFaces(ctx, input)
	if err != nil {
		r.logger.Error("Erro ao comparar faces", zap.Error(err))
		return false, 0, err
	}

	if len(output.FaceMatches) == 0 {
		return false, 0, nil
	}

	return true, *output.FaceMatches[0].Similarity, nil
}

func (r *rekognitionService) CompareWithPreviousPhotos(ctx context.Context, photo core.Photo) (bool, float32, error) {
	// Decodificar a imagem atual de base64
	currentImage, err := base64.StdEncoding.DecodeString(photo.Photo)
	if err != nil {
		r.logger.Error("Erro ao decodificar imagem atual", zap.Error(err))
		return false, 0, err
	}

	// Aqui seria necessário buscar fotos anteriores do mesmo device_id
	// Esta parte será implementada no serviço de telemetria
	return false, 0, errors.New("método não implementado completamente")
}
