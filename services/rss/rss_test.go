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
		var meta models.RSSMeta
		meta.Channel.LastModified = "Thu, 11 Jul 2025 15:04:05 GMT"
		meta.Channel.Title = "Example RSS Feed"
		meta.Channel.Description = "This is a description of the RSS feed."
		_, err := rs.Create(context.Background(), id, "link", meta)
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
		newlink := "newnamewhothat"
		updatedfeed, err := rs.UpdateField(context.Background(), id, "link", newlink)

		require.NoError(t, err)

		if updatedfeed.Link != newlink {
			t.Errorf("Expected updated feed link to be '%s', but got: %s", newlink, updatedfeed.Link)
		}
	})

	t.Run("delete rss", func(t *testing.T) {
		n, err := rs.DeleteByID(context.Background(), id)
		require.NoError(t, err)
		if n != 1 {
			t.Error("expected to delete one entry from db")
		}
	})
}
