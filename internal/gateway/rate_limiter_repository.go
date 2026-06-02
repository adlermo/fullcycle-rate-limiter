package gateway

import (
	"context"
	"time"
)

type RateLimiterRepository interface {
	IsBlocked(
		ctx context.Context,
		key string,
	) (bool, error)

	Increment(
		ctx context.Context,
		key string,
		window time.Duration,
	) (int64, error)

	Block(
		ctx context.Context,
		key string,
		duration time.Duration,
	) error
}
