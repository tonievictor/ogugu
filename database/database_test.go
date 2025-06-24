package database

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/tonievictor/dotenv"
)

func TestDB(t *testing.T) {
	t.Run("incorrect db credentials", func(t *testing.T) {
		_, err := New("pgx", "randomstring")
		require.Error(t, err)
	})

	dotenv.Config("../.env")

	t.Run("correct db credentials", func(t *testing.T) {
		_, err := New("pgx", os.Getenv("DATABASE_URL"))
		require.NoError(t, err)
	})
}
