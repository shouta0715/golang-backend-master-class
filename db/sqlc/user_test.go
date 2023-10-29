package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shouta0715/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	// 作成した時間が同じかどうかを確認する
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, newFullName, updateUser.FullName)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.WithinDuration(t, oldUser.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.PasswordChangedAt, updateUser.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, newEmail, updateUser.Email)
	require.WithinDuration(t, oldUser.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.PasswordChangedAt, updateUser.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)

	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, hashedPassword, updateUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.WithinDuration(t, oldUser.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.PasswordChangedAt, updateUser.PasswordChangedAt, time.Second)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)

	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updateUser)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, hashedPassword, updateUser.HashedPassword)
	require.Equal(t, newFullName, updateUser.FullName)
	require.Equal(t, newEmail, updateUser.Email)
	require.WithinDuration(t, oldUser.CreatedAt, updateUser.CreatedAt, time.Second)
	require.WithinDuration(t, oldUser.PasswordChangedAt, updateUser.PasswordChangedAt, time.Second)
}
