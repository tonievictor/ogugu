package users

import (
	"context"
	"testing"

	"ogugu/services"
)

func TestUserService(t *testing.T) {
	db, tearDownFunc := services.SetupTestDB(t)
	defer tearDownFunc()

	us := New(db) // us -> user service
	id := "uniqueidhaha"

	t.Run("create user", func(t *testing.T) {
		_, err := us.CreateUser(context.Background(), "testusername", "test@random.username", id, "randomavatar")
		if err != nil {
			t.Errorf("Expected no error when creating a user with unique info got %v", err.Error())
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		_, err := us.GetUserByID(context.Background(), id)
		if err != nil {
			t.Errorf("Expected no error when finding user with valid ID, but got: %v", err.Error())
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		_, err := us.GetUserByID(context.Background(), "invalid id")
		if err == nil {
			t.Error("Expected an error when finding user with valid ID, but got none")
		}
	})

	t.Run("update user", func(t *testing.T) {
		newname := "newupdatedtestusername"
		updated, err := us.UpdateUser(context.Background(), id, "username", newname)
		if err != nil {
			t.Errorf("Expected no error when updating an allowed field with the valid details but got %v", err.Error())
		}

		if updated.Username != newname {
			t.Errorf("Expected the row to contain an updated username of %v but got %v", newname, updated.Username)
		}
	})

	t.Run("delete user", func(t *testing.T) {
		err := us.DeleteUserByID(context.Background(), id)
		if err != nil {
			t.Errorf("Expected no error when deleting a valid row with a valid id, but got: %v", err.Error())
		}
	})

	// func (u *UserService) UpdateUser(ctx context.Context, id string, field, value string) (models.User, error) {
}
