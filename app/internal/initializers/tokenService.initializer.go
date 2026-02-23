package initializers

import (
	"JWT/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var TokenService = NewTokenService(os.Getenv("REDIS_ADDRESS"))
var AccessSecret = []byte(os.Getenv("ACCESS_SECRET"))
var RefreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

type RedisTokenService struct {
	redis *redis.Client
}

func NewTokenService(address string) *RedisTokenService {
	return &RedisTokenService{
		redis: redis.NewClient(&redis.Options{
			Addr: address,
		}),
	}
}

func (s *RedisTokenService) StoreRefreshToken(
	ctx context.Context,
	token string,
	userInfo *models.UserInfo,
	duration time.Duration) error {

	jsonUserInfo, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, token, jsonUserInfo, duration).Err()
}

func (s *RedisTokenService) ValidateRefreshToken(ctx context.Context, token string) (userInfo *models.UserInfo, err error) {
	result := s.redis.Get(ctx, token)
	err = result.Err()

	// check redis error
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("token not found or expired")
	} else if err != nil {
		return nil, err
	}

	// unmarshal result
	if err = json.Unmarshal([]byte(result.Val()), &userInfo); err != nil {
		return nil, err
	}

	return
}

func (s *RedisTokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	return s.redis.Del(ctx, token).Err()
}
