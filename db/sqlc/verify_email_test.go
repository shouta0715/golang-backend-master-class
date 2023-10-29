package db

import (
	"context"
	"testing"

	"github.com/shouta0715/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomVerifyEmail(t *testing.T) VerifyEmail {
	user := createRandomUser(t)
	arg := CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      util.RandomEmail(),
		SecretCode: util.RandomString(6),
	}

	verifyEmail, err := testStore.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, arg.Username, verifyEmail.Username)
	require.Equal(t, arg.Email, verifyEmail.Email)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)

	return verifyEmail
}

func TestCreateVerifyEmail(t *testing.T) {
	createRandomVerifyEmail(t)
}
