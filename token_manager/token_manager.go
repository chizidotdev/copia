package token_manager

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type TokenManager interface {
	// CreateToken creates a new token with the given email and duration
	CreateToken(email string, duration time.Duration) (string, error)
	// VerifyToken verifies the given token and returns the payload
	VerifyToken(token string) (*Payload, error)
}
