package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Nickeymaths/bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(10),
		Email:          util.RandomEmail(),
	}
	user, err := testQuery.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	// require.Zero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQuery.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateHashedPasswordOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQuery.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.PasswordChangedAt, oldUser.PasswordChangedAt)
}

func TestUpdateFullNameOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newFullname := util.RandomOwner()

	updatedUser, err := testQuery.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.FullName, newFullname)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.PasswordChangedAt, oldUser.PasswordChangedAt)
}

func TestUpdateEmailOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)
	newEmail := util.RandomEmail()

	updatedUser, err := testQuery.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.PasswordChangedAt, oldUser.PasswordChangedAt)
}

func TestUpdatePasswordChangedAtOnly(t *testing.T) {
	oldUser := CreateRandomUser(t)
	now := time.Now().UTC()

	updatedUser, err := testQuery.UpdateUser(context.Background(), UpdateUserParams{
		PasswordChangedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.NotEqual(t, updatedUser.PasswordChangedAt, oldUser.PasswordChangedAt)
	require.WithinDuration(t, updatedUser.PasswordChangedAt, now, time.Second)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateAllFields(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	now := time.Now().UTC()
	require.NoError(t, err)

	updatedUser, err := testQuery.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		Username: oldUser.Username,
	})

	require.NoError(t, err)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.NotEqual(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.NotEqual(t, updatedUser.FullName, oldUser.FullName)
	require.Equal(t, updatedUser.Email, newEmail)
	require.NotEqual(t, updatedUser.Email, oldUser.Email)
	require.WithinDuration(t, updatedUser.PasswordChangedAt, now, time.Second)
	require.NotEqual(t, updatedUser.PasswordChangedAt, oldUser.PasswordChangedAt)
}
