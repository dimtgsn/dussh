package redis

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

// TODO add fingerprint (?)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	SetRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, token string) error
	UpdateRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetRefreshTokenByUserID(ctx context.Context, userID string) (string, error)
}

type redisCache struct {
	client *redis.Client
}

func MustNew(ctx context.Context, addr, password string) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return &redisCache{client: client}
}

func (rc *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (rc *redisCache) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	err := rc.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *redisCache) SetRefreshToken(
	ctx context.Context,
	userID string,
	token string,
	ttl time.Duration,
) error {
	err := rc.client.HSet(ctx, "refresh_tokens", userID, token).Err()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("refresh_tokens:%s", userID)
	if err := rc.client.Expire(ctx, key, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (rc *redisCache) UpdateRefreshToken(
	ctx context.Context,
	userID string,
	token string,
	ttl time.Duration,
) error {
	err := rc.client.HDel(ctx, "refresh_tokens", userID, token).Err()
	if err != nil {
		return err
	}

	if err := rc.SetRefreshToken(ctx, userID, token, ttl); err != nil {
		return err
	}

	return nil
}

func (rc *redisCache) DeleteRefreshToken(
	ctx context.Context,
	userID string,
	token string,
) error {
	err := rc.client.HDel(ctx, "refresh_tokens", userID, token).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rc *redisCache) GetRefreshTokenByUserID(ctx context.Context, userID string) (string, error) {
	value, err := rc.client.HGet(ctx, "refresh_tokens", userID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("refresh token not found")
		}
		return "", err
	}

	return value, nil
}
