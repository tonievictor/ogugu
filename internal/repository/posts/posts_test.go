package posts

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"ogugu/internal/models"
	"ogugu/internal/repository"
	"ogugu/internal/repository/rss"
)

func TestPostService(t *testing.T) {
	db, teardown := repository.SetupTestDB(t)
	t.Cleanup(teardown)

	ps := New(db)
	rs := rss.New(db)

	rss_id := "rss_id"
	id := "rss_id"

	t.Run("test create rss", func(t *testing.T) {
		var meta models.RSSMeta
		meta.Channel.LastModified = "Thu, 11 Jul 2025 15:04:05 GMT"
		meta.Channel.Title = "Example RSS Feed"
		meta.Channel.Description = "This is a description of the RSS feed."

		_, err := rs.Create(context.Background(), rss_id, "link", meta)
		require.NoError(t, err)
	})

	t.Run("test post creation", func(t *testing.T) {
		p := models.CreatePost{Title: "new", Description: "actually", Link: "www.whocares.com", PubDate: time.Now().Format(time.RFC1123)}
		_, err := ps.CreatePost(context.Background(), id, rss_id, p)
		require.NoError(t, err)
	})

	t.Run("get post by id", func(t *testing.T) {
		_, err := ps.GetByID(context.Background(), id)
		require.NoError(t, err)
	})

	t.Run("fetch all posts", func(t *testing.T) {
		p, err := ps.Fetch(context.Background())
		require.NoError(t, err)

		if len(p) != 1 {
			t.Error("expected one post in the post slice")
		}
	})

	t.Run("delete post by id", func(t *testing.T) {
		n, err := ps.DeletePost(context.Background(), id)
		require.NoError(t, err)
		if n != 1 {
			t.Error("expected to delete one entry from db")
		}
	})
}
