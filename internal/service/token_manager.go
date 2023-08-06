package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chizidotdev/copia/internal/dto"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	Parse(accessToken string) (*dto.Claims, error)
}

type tokenMangerService struct {
	*dto.Claims
	signingKey string
}

func NewTokenManagerService(signingKey string) TokenManager {
	return &tokenMangerService{
		signingKey: signingKey,
		Claims:     &dto.Claims{},
	}
}

func (s *tokenMangerService) Parse(accessToken string) (*dto.Claims, error) {
	splitToken := strings.Split(accessToken, "Bearer ")
	accessToken = splitToken[1]

	token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok { // Check if the signing method is HMAC
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signingKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &dto.Claims{}, errors.New("invalid token")
	}

	if claims["email"] == "" {
		return &dto.Claims{}, errors.New("invalid user email")
	}

	return &dto.Claims{Email: claims["email"].(string)}, nil
}
