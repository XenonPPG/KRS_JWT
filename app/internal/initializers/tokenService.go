package initializers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenService struct {
	redis *redis.Client
}

// TODO: replace with config?
func NewTokenService(address string) *RedisTokenService {
	return &RedisTokenService{
		redis: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (s *RedisTokenService) StoreRefreshToken(ctx context.Context, userID int64, token string, duration time.Duration) error {
	return s.redis.Set(ctx, token, userID, duration).Err()
}

func (s *RedisTokenService) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	userID, err := s.redis.Get(ctx, token).Int64()

	if errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("token not found or expired")
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *RedisTokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.redis.Del(ctx, token).Err()
}
