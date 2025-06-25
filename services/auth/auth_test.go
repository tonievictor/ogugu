package auth

import (
	"context"
	"github.com/stretchr/testify/require"
	"ogugu/services"
	"testing"
)

func TestAuthService(t *testing.T) {
	db, tearDownFunc := services.SetupTestDB(t)
	defer tearDownFunc()

	as := New(db)

	t.Run("test create auth", func(t *testing.T) {
		err := as.CreateAuth(context.Background(), "dummy id", "password")
		require.Error(t, err)
	})

	t.Run("get password with user id", func(t *testing.T) {
		_, err := as.GetPasswordWithUserID(context.Background(), "dummy id")
		require.Error(t, err)
	})
}
