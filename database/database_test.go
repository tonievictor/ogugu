package database

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/tonievictor/dotenv"
)

func TestDB(t *testing.T) {
	t.Run("incorrect db credentials", func(t *testing.T) {
		_, err := New("postgres", "randomstring")
		require.Error(t, err)
	})

	dotenv.Config("../.env")

	t.Run("correct db credentials", func(t *testing.T) {
		_, err := New("postgres", os.Getenv("DATABASE_URL"))
		require.NoError(t, err)
	})
}
