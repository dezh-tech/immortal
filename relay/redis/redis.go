package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client       *redis.Client
	BloomName    string
	QueryTimeout time.Duration
}

func New(cfg Config) (*Redis, error) {
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
		Client:       rc,
		BloomName:    cfg.BloomName,
		QueryTimeout: time.Duration(cfg.QueryTimeout) * time.Millisecond,
	}, nil
}

// ! note: delayed jobs probably are not concurrent safe at the moment.
func (r Redis) AddDelayedJob(listName string,
	data string, delay time.Duration,
) error {
	taskReadyInSeconds := time.Now().Add(delay).Unix()
	member := redis.Z{
		Score:  float64(taskReadyInSeconds),
		Member: data,
	}
	_, err := r.Client.ZAdd(context.Background(), listName, member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r Redis) GetReadyJobs(listName string) ([]string, error) {
	maxTime := time.Now().Unix()

	opt := &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", 0),
		Max: fmt.Sprintf("%d", maxTime),
	}

	cmd := r.Client.ZRevRangeByScore(context.Background(), listName, opt)
	resultSet, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	if err := r.RemoveJobs(listName, resultSet); err != nil {
		return nil, err
	}

	return resultSet, nil
}

func (r Redis) RemoveJobs(listName string, jobs []string) error {
	if len(jobs) == 0 {
		return nil
	}

	_, err := r.Client.ZRem(context.Background(),
		listName, jobs).Result()
	if err != nil {
		return err
	}

	return nil
}
