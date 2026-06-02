package main

import (
	"log"
	"net/http"

	"fullcycle-rate-limiter/internal/config"
	repo "fullcycle-rate-limiter/internal/infra/redis"
	"fullcycle-rate-limiter/internal/middleware"
	"fullcycle-rate-limiter/internal/usecase"

	"github.com/redis/go-redis/v9"
)

func main() {

	cfg := config.Load()

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       0,
		},
	)

	repository := repo.NewRateLimiterRepository(
		redisClient,
	)

	limiter := usecase.NewRateLimiter(
		repository,
		cfg.IpLimit,
		cfg.TokenLimits,
		cfg.BlockDuration,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Write([]byte("ok"))
	})

	handler := middleware.RateLimiter(
		limiter,
	)(mux)

	log.Fatal(
		http.ListenAndServe(
			":"+cfg.Port,
			handler,
		),
	)
}
