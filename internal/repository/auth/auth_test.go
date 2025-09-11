package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"ogugu/internal/repository"
)

func TestAuthService(t *testing.T) {
	db, teardown := repository.SetupTestDB(t)
	t.Cleanup(teardown)

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
