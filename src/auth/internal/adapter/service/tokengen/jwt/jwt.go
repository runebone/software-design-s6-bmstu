package jwt

import (
	"context"
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secretKey string, accessTokenTTLSec, refreshTokenTTLSec time.Duration) *JWTService {
	return &JWTService{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTLSec * 10e9,
		refreshTokenTTL: refreshTokenTTLSec * 10e9,
	}
}

func (s *JWTService) GenerateAccessToken(ctx context.Context, userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.accessTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(s.refreshTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			return "", errors.New("invalid user ID")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return "", errors.New("token expired")
			}
		}

		return userID.String(), nil
	}

	return "", errors.New("invalid token")
}
