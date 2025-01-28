package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"ogugu/services"
)

func TestUserService(t *testing.T) {
	db, tearDownFunc := services.SetupTestDB(t)
	defer tearDownFunc()

	us := New(db) // us -> user service
	id := "uniqueidhaha"

	t.Run("create user", func(t *testing.T) {
		_, err := us.CreateUser(context.Background(), "testusername", "test@random.username", id, "randomavatar")
		require.NoError(t, err)
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
		err := us.DeleteUserByID(context.Background(), id)
		require.NoError(t, err)
	})
}
