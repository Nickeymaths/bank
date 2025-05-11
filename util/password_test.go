package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	// Ensure encode-decode correct
	password := RandomString(6)
	hashedPassword1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	require.NoError(t, IsCorrectPassword(password, hashedPassword1))

	// Ensure 1 password generate different hashs at 2 attempts
	hashedPassword2, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NoError(t, IsCorrectPassword(password, hashedPassword2))
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
