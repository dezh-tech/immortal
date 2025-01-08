package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client              *redis.Client
	BloomFilterName     string
	WhiteListFilterName string
	BlackListFilterName string
	Name                string
	QueryTimeout        time.Duration
}

func New(cfg Config) (*Redis, error) {
	logger.Info("connecting to redis")

	opts, err := redis.ParseURL(cfg.URI)
	if err != nil {
		return nil, err
	}

	rc := redis.NewClient(opts)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnectionTimeout)*time.Millisecond)
	defer cancel()

	// Test the connection
	if err := rc.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return &Redis{
		Client:              rc,
		BloomFilterName:     cfg.BloomFilterName,
		WhiteListFilterName: cfg.WhiteListFilterName,
		BlackListFilterName: cfg.BlackListFilterName,
		QueryTimeout:        time.Duration(cfg.QueryTimeout) * time.Millisecond,
	}, nil
}

func (r *Redis) Close() error {
	logger.Info("closing redis connection")

	return r.Client.Close()
}

// ! note: delayed tasks probably are not concurrent safe at the moment.
func (r *Redis) AddDelayedTask(listName string,
	data string, delay time.Duration,
) error {
	taskReadyInSeconds := time.Now().Add(delay).Unix()
	member := redis.Z{
		Score:  float64(taskReadyInSeconds),
		Member: data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout)
	defer cancel()

	_, err := r.Client.ZAdd(ctx, listName, member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetReadyTasks(listName string) ([]string, error) {
	maxTime := time.Now().Unix()

	opt := &redis.ZRangeBy{
		Min:   fmt.Sprintf("%d", 0),
		Max:   fmt.Sprintf("%d", maxTime),
		Count: 100,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout)
	defer cancel()

	cmd := r.Client.ZRevRangeByScore(ctx, listName, opt)
	resultSet, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	if err := r.RemoveTasks(listName, resultSet); err != nil {
		return nil, err
	}

	return resultSet, nil
}

func (r *Redis) RemoveTasks(listName string, tasks []string) error {
	if len(tasks) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout)
	defer cancel()

	_, err := r.Client.ZRem(ctx,
		listName, tasks).Result()
	if err != nil {
		return err
	}

	return nil
}
