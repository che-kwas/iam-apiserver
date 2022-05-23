// Package redis defines the global redis client.
package redis

import (
	"sync"

	"github.com/go-redis/redis/v8"

	"iam-apiserver/internal/pkg/config"
)

var (
	rdb  redis.UniversalClient
	once sync.Once
)

// Client returns the global redis client.
func Client() redis.UniversalClient {
	if rdb != nil {
		return rdb
	}

	once.Do(func() {
		opts := config.Cfg().RedisOpts
		rdb = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    opts.Addrs,
			Password: opts.Password,
			DB:       opts.Database,
		})
	})

	return rdb
}
