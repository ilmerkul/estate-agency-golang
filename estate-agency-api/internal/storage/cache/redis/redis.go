package redis

import (
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func New() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "3306",
			"server2": "3306",
		},
	})

	cacheClient := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return cacheClient
}
