package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	context    context.Context
	client     *redis.Client
	expiration time.Duration
	logger     *slog.Logger
}

func New(ctx context.Context, address string, password string, logger *slog.Logger, expiration time.Duration) Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})

	return Storage{
		context:    ctx,
		client:     client,
		expiration: expiration,
		logger:     logger,
	}
}

func (s Storage) Save(key string, value interface{}) error {
	err := s.client.Set(s.context, key, value, s.expiration).Err()
	if err != nil {
		s.logger.Error(fmt.Sprintf("redis error: %s", err))
	}

	return err
}

func (s Storage) Read(key string) (interface{}, error) {
	val, err := s.client.Get(s.context, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s Storage) ReadAllKeys() ([]string, error) {
	rkeys := s.client.Keys(s.context, "*")

	return rkeys.Result()
}
