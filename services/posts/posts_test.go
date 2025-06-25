package posts

import (
	"context"
	"github.com/stretchr/testify/require"
	"ogugu/models"
	"ogugu/services"
	"ogugu/services/rss"
	"testing"
	"time"
)

func TestPostService(t *testing.T) {
	db, teardown := services.SetupTestDB(t)
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
}
