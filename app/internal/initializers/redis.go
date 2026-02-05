package initializers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenService struct {
	redis *redis.Client
}

// TODO: replace with config?
func (s *TokenService) NewTokenService(address string) {
	s.redis = redis.NewClient(&redis.Options{
		Addr: address,
	})
}

func (s *TokenService) StoreRefreshToken(ctx context.Context, userID int64, token string, duration time.Duration) error {
	return s.redis.Set(ctx, token, userID, duration).Err()
}

func (s *TokenService) ValidateRefreshToken(ctx context.Context, token string) (int64, error) {
	userID, err := s.redis.Get(ctx, token).Int64()

	if errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("token not found or expired")
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.redis.Del(ctx, token).Err()
}
