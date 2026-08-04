package handler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type dingTalkRedisAppTokenStore struct {
	rdb redis.Cmdable
}

func NewDingTalkRedisAppTokenStore(rdb redis.Cmdable) dingTalkAppTokenStore {
	if rdb == nil {
		return nil
	}
	return &dingTalkRedisAppTokenStore{rdb: rdb}
}

func (s *dingTalkRedisAppTokenStore) Get(ctx context.Context, key string) (string, error) {
	return s.rdb.Get(ctx, key).Result()
}

func (s *dingTalkRedisAppTokenStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.rdb.Set(ctx, key, value, ttl).Err()
}
