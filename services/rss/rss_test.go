package rss

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"ogugu/services"
)

func TestRssService(t *testing.T) {
	db, tearDownFunc := services.SetupTestDB(t)
	defer tearDownFunc()

	rs := New(db)
	id := "uniqueidhaha"

	t.Run("test create rss", func(t *testing.T) {
		_, err := rs.Create(context.Background(), "link", "name", id)
		require.NoError(t, err)
	})

	t.Run("test create rss", func(t *testing.T) {
		_, err := rs.Create(context.Background(), "link", "name", id)
		require.Error(t, err)
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
