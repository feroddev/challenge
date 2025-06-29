package recognition

import (
	"bytes"

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

func (c *AWSRekognitionClient) CompareFaces(sourceImage []byte, targetImage []byte) (float32, error) {
	input := &rekognition.CompareFacesInput{
		SourceImage: &rekognition.Image{
			Bytes: sourceImage,
		},
		TargetImage: &rekognition.Image{
			Bytes: targetImage,
		},
		SimilarityThreshold: aws.Float64(0),
	}

	result, err := c.client.CompareFaces(input)
	if err != nil {
		c.logger.Error("Erro ao comparar faces", zap.Error(err))
		return 0, err
	}

	if len(result.FaceMatches) == 0 {
		return 0, nil
	}

	var maxSimilarity float32
	for _, match := range result.FaceMatches {
		similarity := float32(*match.Similarity)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
		}
	}

	return maxSimilarity / 100.0, nil
}
