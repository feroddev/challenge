package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github/feroddev/challengeV3/internal/core"
	"github/feroddev/challengeV3/internal/repositories"
)

type AuthServiceConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
}

type AuthService struct {
	config AuthServiceConfig
	repo   *repositories.UserRepository
	logger *zap.Logger
}

func NewAuthService(config AuthServiceConfig, db *gorm.DB, logger *zap.Logger) *AuthService {
	return &AuthService{
		config: config,
		repo:   repositories.NewUserRepository(db),
		logger: logger,
	}
}

func (s *AuthService) Register(user core.User) error {
	existingUser, err := s.repo.FindByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("usuario ja existe")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	s.logger.Info("Registrando novo usuario", 
		zap.String("username", user.Username), 
		zap.Int("password_length", len(user.Password)))

	return s.repo.Save(&user)
}

func (s *AuthService) Login(request core.LoginRequest) (*core.TokenResponse, error) {
	user, err := s.repo.FindByUsername(request.Username)
	if err != nil {
		s.logger.Error("Erro ao buscar usuario", zap.Error(err))
		return nil, errors.New("credenciais invalidas")
	}

	s.logger.Info("Login bem-sucedido", 
		zap.String("username", request.Username))

	expirationTime := time.Now().Add(s.config.TokenDuration)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &core.TokenResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}, nil
}

func (s *AuthService) GetUserByID(id uint) (*core.User, error) {
	return s.repo.FindByID(id)
}

func (s *AuthService) ValidateToken(tokenString string) (*core.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo de assinatura invalido")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		username := claims["username"].(string)
		role := claims["role"].(string)

		return &core.Claims{
			UserID:   userID,
			Username: username,
			Role:     role,
		}, nil
	}

	return nil, errors.New("token invalido")
}
