package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/storage"
)

type TokenStorage struct {
	db *redis.Client
}

func New(addr, password string) *TokenStorage {
	return &TokenStorage{
		db: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		}),
	}
}

func (db *TokenStorage) Close() {
	db.db.Close()
}

func (db *TokenStorage) JWT(ctx context.Context, telegramLogin string) (bool, error) {
	const op = "storage.redis.JWT"

	key := fmt.Sprintf("telegram_login:%s", telegramLogin)

	_, err := db.db.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, fmt.Errorf("%s: %w", op, storage.ErrTokenNotFound)
		}
		return false, fmt.Errorf("%s: failed to get JWT for user %s: %w", op, telegramLogin, err)
	}

	return true, nil
}

func (db *TokenStorage) SaveJWT(ctx context.Context, token string, telegramLogin string, ttl time.Duration) error {
	const op = "storage.redis.SaveJWT"

	key := fmt.Sprintf("telegram_login:%s", telegramLogin)

	wasSet, err := db.db.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !wasSet {
		return fmt.Errorf("%s %w", op, storage.ErrTokenExists)
	}

	return nil
}

func (db *TokenStorage) DeleteJWT(ctx context.Context, telegramLogin string) error {
	const op = "storage.redis.DeleteJWT"

	key := fmt.Sprintf("telegram_login:%s", telegramLogin)

	err := db.db.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
