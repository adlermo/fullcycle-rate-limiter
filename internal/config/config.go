package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	RedisAddr     string
	RedisPassword string

	IpLimit int

	TokenLimits map[string]int

	BlockDuration time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	ipLimit, _ := strconv.Atoi(
		getEnv("RATE_LIMIT_IP", "10"),
	)

	blockSeconds, _ := strconv.Atoi(
		getEnv("BLOCK_DURATION_SECONDS", "300"),
	)

	tokenLimits := make(map[string]int)

	for _, env := range os.Environ() {

		parts := strings.SplitN(env, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		if !strings.HasPrefix(key, "TOKEN_") {
			continue
		}

		limit, err := strconv.Atoi(value)

		if err != nil {
			continue
		}

		token := strings.TrimPrefix(
			key,
			"TOKEN_",
		)

		tokenLimits[token] = limit
	}

	return &Config{
		Port: getEnv(
			"PORT",
			"8080",
		),

		RedisAddr: getEnv(
			"REDIS_ADDR",
			"redis:6379",
		),

		RedisPassword: getEnv(
			"REDIS_PASSWORD",
			"",
		),

		IpLimit: ipLimit,

		TokenLimits: tokenLimits,

		BlockDuration: time.Duration(
			blockSeconds,
		) * time.Second,
	}
}

func getEnv(
	key string,
	defaultValue string,
) string {

	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
