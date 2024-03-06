package token

import (
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
	"github.com/yuuban007/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
	// test short key
	maker, err := NewPasetoMaker(util.RandomString(11))
	require.Error(t, err, ErrSecretKeyTooShort)
	require.Empty(t, maker)

	// test happy case
	maker, err = NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()

	token, err := maker.CreateToken(username, -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {

	// create a token with different symmetricKey
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	wrongKey := util.RandomString(32)
	token, err := paseto.NewV2().Encrypt([]byte(wrongKey), payload, nil)
	require.NoError(t, err)

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}
