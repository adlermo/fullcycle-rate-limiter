package repository

import (
	"context"
	"fullcycle-rate-limiter/internal/gateway"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterRepository struct {
	client *redis.Client
}

func NewRateLimiterRepository(
	client *redis.Client,
) *RateLimiterRepository {
	return &RateLimiterRepository{
		client: client,
	}
}

func (r *RateLimiterRepository) IsBlocked(
	ctx context.Context,
	key string,
) (bool, error) {

	result, err := r.client.Exists(
		ctx,
		"blocked:"+key,
	).Result()

	if err != nil {
		return false, err
	}

	return result > 0, nil
}

func (r *RateLimiterRepository) Block(
	ctx context.Context,
	key string,
	duration time.Duration,
) error {

	return r.client.Set(
		ctx,
		"blocked:"+key,
		1,
		duration,
	).Err()
}

func (r *RateLimiterRepository) Increment(
	ctx context.Context,
	key string,
	window time.Duration,
) (int64, error) {

	count, err := r.client.Incr(
		ctx,
		"counter:"+key,
	).Result()

	if err != nil {
		return 0, err
	}

	if count == 1 {
		err = r.client.Expire(
			ctx,
			"counter:"+key,
			window,
		).Err()

		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

var _ gateway.RateLimiterRepository = (*RateLimiterRepository)(nil)
