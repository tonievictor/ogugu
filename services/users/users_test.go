package users

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"ogugu/models"
	"ogugu/services"
)

func TestUserService(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	mfile := "file://" + filepath.Dir(filepath.Dir(dir)) + "/migrations"
	db, teardown := services.SetupTestDB(t, mfile)
	defer teardown()

	us := New(db) // us -> user service
	id := "uniqueidhaha"

	t.Run("create user", func(t *testing.T) {
		body := models.CreateUserBody{
			Username: "testusername",
			Email:    "test@random.username",
			Avatar:   "randomavatar",
			Password: "password",
		}
		_, err := us.CreateUser(context.Background(), id, body)
		require.NoError(t, err)
	})

	t.Run("create user", func(t *testing.T) {
		body := models.CreateUserBody{
			Username: "testusername",
			Email:    "test@random.username",
			Avatar:   "randomavatar",
			Password: "password",
		}
		_, err := us.CreateUser(context.Background(), id, body)
		require.Error(t, err)
	})

	t.Run("get user with auth", func(t *testing.T) {
		_, _, err := us.GetUserAuth(context.Background(), "test@random.username")
		require.Error(t, err)
	})

	t.Run("get user by id", func(t *testing.T) {
		_, err := us.GetUserByID(context.Background(), id)
		require.NoError(t, err)
	})

	t.Run("get user by id", func(t *testing.T) {
		_, err := us.GetUserByID(context.Background(), "invalid id")
		require.Error(t, err)
	})

	t.Run("update user", func(t *testing.T) {
		newname := "newupdatedtestusername"
		updated, err := us.UpdateUser(context.Background(), id, "username", newname)
		require.NoError(t, err)

		if updated.Username != newname {
			t.Errorf("Expected the row to contain an updated username of %v but got %v", newname, updated.Username)
		}
	})

	t.Run("delete user", func(t *testing.T) {
		n, err := us.DeleteUserByID(context.Background(), id)
		require.NoError(t, err)

		if n != 1 {
			t.Error("expected to delete one entry from db")
		}
	})
}
