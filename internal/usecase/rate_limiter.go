package usecase

import (
	"context"
	"time"

	"fullcycle-rate-limiter/internal/gateway"
)

type RateLimiter struct {
	repo gateway.RateLimiterRepository

	ipLimit       int
	tokenLimits   map[string]int
	blockDuration time.Duration
}

func NewRateLimiter(
	repo gateway.RateLimiterRepository,
	ipLimit int,
	tokenLimits map[string]int,
	blockDuration time.Duration,
) *RateLimiter {
	return &RateLimiter{
		repo:          repo,
		ipLimit:       ipLimit,
		tokenLimits:   tokenLimits,
		blockDuration: blockDuration,
	}
}

func (r *RateLimiter) Allow(
	ctx context.Context,
	ip string,
	token string,
) (bool, error) {

	var (
		key   string
		limit int
	)

	if token != "" {
		if tokenLimit, ok := r.tokenLimits[token]; ok {
			key = "token:" + token
			limit = tokenLimit
		}
	}

	if key == "" {
		key = "ip:" + ip
		limit = r.ipLimit
	}

	blocked, err := r.repo.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}

	if blocked {
		return false, nil
	}

	count, err := r.repo.Increment(
		ctx,
		key,
		time.Second,
	)
	if err != nil {
		return false, err
	}

	if count > int64(limit) {
		err = r.repo.Block(
			ctx,
			key,
			r.blockDuration,
		)

		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
