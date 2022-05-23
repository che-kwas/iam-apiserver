// Package cache define the global cache instance.
package cache

import (
	"sync"

	"github.com/che-kwas/iam-kit/cache"
	"github.com/go-redis/redis/v8"

	"iam-apiserver/internal/apiserver/config"
)

var (
	rdb  *redis.UniversalClient
	once sync.Once
)

// Cache returns the global config instance.
func Cache() *redis.UniversalClient {
	if rdb != nil {
		return rdb
	}

	once.Do(func() {
		rdb, _ = cache.NewRedisIns(config.Cfg().RedisOpts)
	})

	return rdb
}
