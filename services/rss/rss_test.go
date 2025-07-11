package rss

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"ogugu/models"
	"ogugu/services"
)

func TestRssService(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	mfile := "file://" + filepath.Dir(filepath.Dir(dir)) + "/migrations"
	db, teardown := services.SetupTestDB(t, mfile)
	defer teardown()

	rs := New(db)
	id := "uniqueidhaha"

	t.Run("test create rss", func(t *testing.T) {
		body := models.CreateRssBody{Link: "link", Name: "name"}
		_, err := rs.Create(context.Background(), id, body)
		require.NoError(t, err)
	})

	t.Run("test find rss by id", func(t *testing.T) {
		_, err := rs.FindByID(context.Background(), id)
		require.NoError(t, err)
	})

	t.Run("test find rss by id", func(t *testing.T) {
		_, err := rs.FindByID(context.Background(), "non-existent")
		require.Error(t, err)
	})

	t.Run("test fetch all rss", func(t *testing.T) {
		_, err := rs.Fetch(context.Background())
		require.NoError(t, err)
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		updatedfeed, err := rs.Update(context.Background(), id, "name", newname)

		require.NoError(t, err)

		if updatedfeed.Name != newname {
			t.Errorf("Expected updated feed name to be '%s', but got: %s", newname, updatedfeed.Name)
		}
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		_, err := rs.Update(context.Background(), id, "id", newname)
		require.Error(t, err)
	})

	t.Run("update rss", func(t *testing.T) {
		newname := "newnamewhothat"
		_, err := rs.Update(context.Background(), "nonexistentid", "id", newname)
		require.Error(t, err)
	})

	t.Run("delete rss", func(t *testing.T) {
		err := rs.DeleteByID(context.Background(), id)

		require.NoError(t, err)
	})
}
