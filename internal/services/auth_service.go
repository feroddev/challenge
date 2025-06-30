package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github/feroddev/challengeV3/internal/core"
)

type AuthServiceConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
}

type AuthService struct {
	config AuthServiceConfig
	users  []core.User
	logger *zap.Logger
}

func NewAuthService(config AuthServiceConfig, logger *zap.Logger) *AuthService {
	return &AuthService{
		config: config,
		users:  make([]core.User, 0),
		logger: logger,
	}
}

func (s *AuthService) Register(user core.User) error {
	for _, existingUser := range s.users {
		if existingUser.Username == user.Username {
			return errors.New("usuario ja existe")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.ID = uint(len(s.users) + 1)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	s.users = append(s.users, user)
	return nil
}

func (s *AuthService) Login(request core.LoginRequest) (*core.TokenResponse, error) {
	var user *core.User
	for _, u := range s.users {
		if u.Username == request.Username {
			user = &u
			break
		}
	}

	if user == nil {
		return nil, errors.New("credenciais invalidas")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New("credenciais invalidas")
	}

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
