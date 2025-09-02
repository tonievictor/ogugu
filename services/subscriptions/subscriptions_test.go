package subscriptions

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"ogugu/models"
	"ogugu/services"
	"ogugu/services/rss"
	"ogugu/services/users"
)

func TestSubscriptionService(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	mfile := "file://" + filepath.Dir(filepath.Dir(dir)) + "/migrations"
	db, teardown := services.SetupTestDB(t, mfile)
	t.Cleanup(teardown)

	rssid := "rssid"
	userid := "userid"
	subid := "userid"
	us := users.New(db)
	rs := rss.New(db)
	ss := New(db)

	var meta models.RSSMeta
	meta.Channel.LastModified = "Thu, 11 Jul 2025 15:04:05 GMT"
	meta.Channel.Title = "Example RSS Feed"
	meta.Channel.Description = "This is a description of the RSS feed."
	_, err = rs.Create(context.Background(), rssid, "rsslink", meta)
	require.NoError(t, err)

	var createUser models.CreateUserBody
	createUser.Username = "username"
	createUser.Password = "password"
	createUser.Avatar = "avatar"
	createUser.Email = "email"
	_, err = us.CreateUser(context.Background(), userid, createUser)

	t.Run("create subscription", func(t *testing.T) {
		_, err := ss.CreateSub(context.Background(), subid, userid, rssid)
		require.NoError(t, err)
	})

	t.Run("create subscription", func(t *testing.T) {
		_, err := ss.CreateSub(context.Background(), subid, userid, rssid)
		require.Error(t, err)
	})

	t.Run("create subscription", func(t *testing.T) {
		_, err := ss.CreateSub(context.Background(), "id2", userid, rssid)
		require.Error(t, err)
	})

	t.Run("get subscriptions", func(t *testing.T) {
		_, err := ss.GetSubs(context.Background())
		require.NoError(t, err)
	})

	t.Run("get subscriptions by user id", func(t *testing.T) {
		_, err := ss.GetSubsByUserID(context.Background(), userid)
		require.NoError(t, err)
	})

	t.Run("get subscription by id", func(t *testing.T) {
		_, err := ss.GetSubByID(context.Background(), subid)
		require.NoError(t, err)
	})

	t.Run("unsubscribe", func(t *testing.T) {
		n, err := ss.DeleteSub(context.Background(), userid, rssid)
		require.NoError(t, err)

		if n != 1 {
			t.Error("expected to delete one entry from subscriptions table")
		}
	})

	t.Run("get subscriptions from user post", func(t *testing.T) {
		_, err := ss.GetPostFromSubScriptions(context.Background(), userid)
		require.NoError(t, err)
	})
}
