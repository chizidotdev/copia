package token_manager

import (
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoTokenManager struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoTokenManager(symmetricKey []byte) (TokenManager, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, ErrInvalidToken
	}

	return &PasetoTokenManager{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}, nil
}

func (p PasetoTokenManager) CreateToken(email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

// VerifyToken verifies the given token and returns the payload
func (p PasetoTokenManager) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, &payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
