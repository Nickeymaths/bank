package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

func NewPasetoMarker(secretKey string) (Maker, error) {
	if len(secretKey) < minKeyLen {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minKeyLen)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: secretKey,
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := maker.paseto.Encrypt([]byte(maker.symmetricKey), &payload, nil)
	return token, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Decrypt(token, []byte(maker.symmetricKey), &payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
