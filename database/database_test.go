package database

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tonievictor/dotenv"
)

func TestDB(t *testing.T) {
	t.Run("incorrect db credentials", func(t *testing.T) {
		_, err := Setup("postgres", "randomstring")
		if err == nil {
			t.Error("incorrect db credentials should fail")
		}
	})

	dotenv.Config("../.env")

	t.Run("correct db credentials", func(t *testing.T) {
		_, err := Setup("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			t.Error("correct db credentials should pass")
		}
	})
}
