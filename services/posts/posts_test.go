package posts

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"ogugu/models"
	"ogugu/services"
	"ogugu/services/rss"
)

func TestPostService(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	mfile := "file://" + filepath.Dir(filepath.Dir(dir)) + "/migrations"
	db, teardown := services.SetupTestDB(t, mfile)
	defer teardown()

	ps := New(db)
	rs := rss.New(db)

	rss_id := "rss_id"
	id := "rss_id"

	t.Run("test create rss", func(t *testing.T) {
		body := models.CreateRssBody{Link: "link", Name: "name"}
		_, err := rs.Create(context.Background(), rss_id, body)
		require.NoError(t, err)
	})

	t.Run("test post creation", func(t *testing.T) {
		p := models.CreatePost{Title: "new", Description: "actually", Link: "www.whocares.com", PubDate: time.Now()}
		_, err := ps.CreatePost(context.Background(), id, rss_id, p)
		require.NoError(t, err)
	})

	t.Run("get post by id", func(t *testing.T) {
		_, err := ps.GetPostByID(context.Background(), id)
		require.NoError(t, err)
	})

	t.Run("fetch all posts", func(t *testing.T) {
		p, err := ps.FetchPosts(context.Background())
		require.NoError(t, err)

		if len(p) != 1 {
			t.Error("expected one post in the post slice")
		}
	})

	t.Run("delete post by id", func(t *testing.T) {
		err := ps.DeletePost(context.Background(), id)
		require.NoError(t, err)
	})
}
