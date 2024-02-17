package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)
	hashedPassword1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(8)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.Error(t, err)

	hashedPassword2, _ := HashPassword(password)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
	err = CheckPassword(password, hashedPassword2)
	require.NoError(t, err)
}
