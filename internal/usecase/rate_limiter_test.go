package usecase_test

import (
	"context"
	"testing"
	"time"

	"fullcycle-rate-limiter/internal/usecase"
)

func TestShouldAllowRequestBelowLimit(t *testing.T) {
	repo := NewFakeRepository()

	limiter := usecase.NewRateLimiter(
		repo,
		10,
		nil,
		time.Minute,
	)

	allowed, err := limiter.Allow(
		context.Background(),
		"127.0.0.1",
		"",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !allowed {
		t.Fatal("request should be allowed")
	}
}

func TestShouldBlockRequestAboveLimit(t *testing.T) {
	repo := NewFakeRepository()

	limiter := usecase.NewRateLimiter(
		repo,
		2,
		nil,
		time.Minute,
	)

	ctx := context.Background()

	_, _ = limiter.Allow(ctx, "127.0.0.1", "")
	_, _ = limiter.Allow(ctx, "127.0.0.1", "")

	allowed, err := limiter.Allow(ctx, "127.0.0.1", "")

	if err != nil {
		t.Fatal(err)
	}

	if allowed {
		t.Fatal("request should be blocked")
	}

	blocked, _ := repo.IsBlocked(ctx, "ip:127.0.0.1")

	if !blocked {
		t.Fatal("ip should be blocked")
	}
}

func TestTokenShouldOverrideIPLimit(t *testing.T) {
	repo := NewFakeRepository()

	limiter := usecase.NewRateLimiter(
		repo,
		2,
		map[string]int{
			"premium-token": 5,
		},
		time.Minute,
	)

	ctx := context.Background()

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(
			ctx,
			"127.0.0.1",
			"premium-token",
		)

		if err != nil {
			t.Fatal(err)
		}

		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	allowed, err := limiter.Allow(
		ctx,
		"127.0.0.1",
		"premium-token",
	)

	if err != nil {
		t.Fatal(err)
	}

	if allowed {
		t.Fatal("token should be blocked after reaching its own limit")
	}
}
