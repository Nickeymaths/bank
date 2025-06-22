package db

import (
	"context"
	"testing"

	"github.com/Nickeymaths/bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateVerifyEmail(t *testing.T) {
	user := CreateRandomUser(t)
	arg := CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	verifyEmail, err := testQuery.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, verifyEmail.ID)
	require.Equal(t, user.Username, verifyEmail.Username)
	require.Equal(t, user.Email, verifyEmail.Email)
	require.Equal(t, false, verifyEmail.IsUsed)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)
	require.Equal(t, 32, len(verifyEmail.SecretCode))
	require.NotEmpty(t, verifyEmail.CreatedAt)
	require.NotEmpty(t, verifyEmail.ExpiredAt)
}

func TestUpdateVerifyStatus(t *testing.T) {
	user := CreateRandomUser(t)
	arg := CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	verifyEmail, err := testQuery.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)

	updatedVerifyEmail, err := testQuery.UpdateVerifyStatus(context.Background(), UpdateVerifyStatusParams{
		ID:         verifyEmail.ID,
		SecretCode: verifyEmail.SecretCode,
	})
	require.NoError(t, err)
	require.Equal(t, true, updatedVerifyEmail.IsUsed)
}
