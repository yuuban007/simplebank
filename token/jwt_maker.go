package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var minServerSecretKeySize = 32

var (
	ErrSecretKeyTooShort = fmt.Errorf("given secret key should be at least %d characters", minServerSecretKeySize)
)

// JWTMaker is a Json Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker returns a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minServerSecretKeySize {
		return nil, ErrSecretKeyTooShort
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken implements Maker.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	Payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken implements Maker.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if ok {
			return nil, ErrInvalidToken
		}
		// if the convert success, it meaning the alg match
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrorTokenExpired) {
			return nil, ErrorTokenExpired
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
