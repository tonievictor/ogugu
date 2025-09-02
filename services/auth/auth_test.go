package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"ogugu/services"
)

func TestAuthService(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	mfile := "file://" + filepath.Dir(filepath.Dir(dir)) + "/migrations"
	db, teardown := services.SetupTestDB(t, mfile)
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
