package cache

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

//REFRESH-TOKEN

func (r *RedisService) TokenSetWithExpiry(ctx context.Context, userID string, token string, expiry time.Duration) (err error) {
	s := config.Conf.PrefixForRedis.RefreshToken + token
	err = r.client.Set(ctx, s, userID, expiry).Err()
	return
}

func (r *RedisService) GetUserIDByToken(ctx context.Context, token string) (userID *string, err error) {
	s := config.Conf.PrefixForRedis.RefreshToken + token
	usID, err := r.client.Get(ctx, s).Result()
	if err != nil {
		if err == redis.Nil {
			userID = nil
			err = nil
			return
		}
		return
	}
	userID = &usID

	return
}

func (r *RedisService) DeletePreviousRefreshToken(ctx context.Context, token string) (err error) {
	key := config.Conf.PrefixForRedis.RefreshToken + token
	err = r.client.Del(ctx, key).Err()
	return
}
