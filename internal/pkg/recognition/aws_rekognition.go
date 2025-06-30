package recognition

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"go.uber.org/zap"
)

type AWSRekognitionClient struct {
	client *rekognition.Rekognition
	logger *zap.Logger
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func NewAWSRekognitionClient(config AWSConfig, logger *zap.Logger) (*AWSRekognitionClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretAccessKey,
			config.SessionToken,
		),
	})

	if err != nil {
		logger.Error("Erro ao criar sessao AWS", zap.Error(err))
		return nil, err
	}

	client := rekognition.New(sess)

	return &AWSRekognitionClient{
		client: client,
		logger: logger,
	}, nil
}

func (c *AWSRekognitionClient) CompareFaces(sourceImageBase64 string, targetImageBase64 string) (bool, float32, error) {
	sourceBytes, err := base64.StdEncoding.DecodeString(sourceImageBase64)
	if err != nil {
		c.logger.Error("Erro ao decodificar imagem fonte", zap.Error(err))
		return false, 0, err
	}

	targetBytes, err := base64.StdEncoding.DecodeString(targetImageBase64)
	if err != nil {
		c.logger.Error("Erro ao decodificar imagem alvo", zap.Error(err))
		return false, 0, err
	}

	input := &rekognition.CompareFacesInput{
		SourceImage: &rekognition.Image{
			Bytes: sourceBytes,
		},
		TargetImage: &rekognition.Image{
			Bytes: targetBytes,
		},
		SimilarityThreshold: aws.Float64(80.0),
	}

	result, err := c.client.CompareFaces(input)
	if err != nil {
		c.logger.Error("Erro ao comparar faces", zap.Error(err))
		return false, 0, err
	}

	if len(result.FaceMatches) == 0 {
		c.logger.Info("Nenhuma correspondência facial encontrada")
		return false, 0, nil
	}

	similarity := float32(*result.FaceMatches[0].Similarity)
	c.logger.Info("Faces comparadas com sucesso",
		zap.Float32("similarity", similarity))

	return true, similarity, nil
}
