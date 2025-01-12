package rss

import (
	"context"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tonievictor/dotenv"

	"ogugu/database"
)

func TestRssService(t *testing.T) {
	dotenv.Config("../../.env")
	db, err := database.Setup("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Errorf("There was an error setting up the database. Error: %s", err.Error())
		return
	}

	rs := New(db)
	id := "uniqueidhaha"

	t.Run("test create rss", func(t *testing.T) {
		_, err := rs.Create(context.Background(), "link", "name", id)
		if err != nil {
			t.Errorf("Expected no error when creating RSS with correct information, but got: %v", err)
		}
	})

	t.Run("test create rss", func(t *testing.T) {
		_, err := rs.Create(context.Background(), "link", "name", id)
		if err == nil {
			t.Error("Expected error when creating RSS with conflicting information, but no error was returned.")
		}
	})

	t.Run("test find rss by id", func(t *testing.T) {
		_, err := rs.FindByID(context.Background(), id)
		if err != nil {
			t.Errorf("Expected no error when finding RSS with valid ID, but got: %v", err)
		}
	})

	t.Run("test find rss by id", func(t *testing.T) {
		_, err := rs.FindByID(context.Background(), "non-existent")
		if err == nil {
			t.Errorf("Expected error when finding RSS with non-existent ID, but no error was returned")
		}
	})

	t.Run("test fetch all rss", func(t *testing.T) {
		_, err := rs.Fetch(context.Background())
		if err != nil {
			t.Errorf("Expected no error when fetching all RSS, but got: %v", err)
		}
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		updatedfeed, err := rs.Update(context.Background(), id, "name", newname)
		if err != nil {
			t.Errorf("Expected no error when updating a valid field, but got: %v", err)
		}

		if updatedfeed.Name != newname {
			t.Errorf("Expected updated feed name to be '%s', but got: %s", newname, updatedfeed.Name)
		}
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		_, err := rs.Update(context.Background(), id, "id", newname)
		if err == nil {
			t.Error("Expected an error when updating an illegal field, but got none")
		}
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		_, err := rs.Update(context.Background(), "nonexistentid", "id", newname)
		if err == nil {
			t.Error("Expected an error when updating a non existent row, but got none")
		}
	})

	t.Run("delete rss", func(t *testing.T) {
		err := rs.DeleteByID(context.Background(), id)
		if err != nil {
			t.Errorf("Expected no error when deleting a valid row with a valid id, but got: %v", err)
		}
	})
}
