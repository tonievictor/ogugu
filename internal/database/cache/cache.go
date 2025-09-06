package cache

import (
	"github.com/redis/go-redis/v9"
)

func Setup(conn_url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(conn_url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return client, nil
}
