package token

import (
	"testing"
	"time"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

func TestCreatePasetoToken(t *testing.T) {
	username := utils.RandomName()

	duration := time.Minute
	createAt := time.Now()
	expireAt := createAt.Add(duration)

	paseto, err := NewPaseto(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, paseto)

	token, err := paseto.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, paseto)

	payload, err := paseto.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, paseto)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, createAt, payload.CreatedAt, time.Second)
	require.WithinDuration(t, expireAt, payload.ExpiryAt, time.Second)

}

func TestExpiredToken(t *testing.T) {
	paseto, err := NewPaseto(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, paseto)

	token, err := paseto.CreateToken(utils.RandomName(), -time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, paseto)

	payload, err := paseto.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)

}
