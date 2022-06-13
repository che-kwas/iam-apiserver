// Package redis implements `iam-apiserver/internal/apiserver/publisher.Publisher` interface.
package redis

import (
	"context"

	rdb "github.com/che-kwas/iam-kit/redis"
	"github.com/go-redis/redis/v8"

	"iam-apiserver/internal/apiserver/publisher"
)

type redisPub struct {
	cli redis.UniversalClient
}

func (r *redisPub) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.cli.Publish(ctx, channel, message).Err()
}

func (r *redisPub) Close() error {
	return r.cli.Close()
}

// NewRedisPub returns a redis publisher.
func NewRedisPub() (publisher.Publisher, error) {
	cli, err := rdb.NewRedisIns()
	if err != nil {
		return nil, err
	}

	return &redisPub{cli}, nil
}
