package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("faulty redis connection string", func(t *testing.T) {
		_, err := Setup("redis://<user>:<pass>@localhost:6379/<db>")
		if err == nil {
			t.Error("wrong redis url should fail")
		}
	})

	t.Run("correct redis url", func(t *testing.T) {
		_, err := Setup("redis://localhost:6379")
		if err != nil {
			t.Error("correct redis url should pass")
		}
	})
}
