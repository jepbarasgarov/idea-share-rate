package cache

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func (r *RedisService) WorkerRestrictWithExpiry(ctx context.Context, fullname string, expiry time.Duration) (err error) {
	allUpper := strings.ToUpper(fullname)
	s := config.Conf.PrefixForRedis.WorkerRestrictIdeaSubmit + allUpper
	err = r.client.Set(ctx, s, true, expiry).Err()
	return
}

func (r *RedisService) HasRestrctionForWorker(ctx context.Context, fullname string) (hasBlock bool, err error) {
	allUpper := strings.ToUpper(fullname)
	key := config.Conf.PrefixForRedis.WorkerRestrictIdeaSubmit + allUpper
	_, err = r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			hasBlock = false
			err = nil
			return
		}
		return
	}

	hasBlock = true
	return
}
