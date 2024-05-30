package test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"go-rate-limiter/limiter"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	config := limiter.Config{
		RateLimit:     5,
		BlockDuration: 10 * time.Second,
	}

	rateLimiter := limiter.NewRedisRateLimiter(config)

	ctx := context.Background()
	key := "test_key"

	// Limpar o estado anterior do Redis
	client.Del(ctx, key)

	// Testar que até 5 requisições são permitidas
	for i := 0; i < 5; i++ {
		allowed, err := rateLimiter.Allow(ctx, key)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// A 6ª requisição deve ser bloqueada
	allowed, err := rateLimiter.Allow(ctx, key)
	assert.NoError(t, err)
	assert.False(t, allowed)

	// Esperar o tempo de bloqueio
	time.Sleep(config.BlockDuration)

	// Após o tempo de bloqueio, a requisição deve ser permitida novamente
	allowed, err = rateLimiter.Allow(ctx, key)
	assert.NoError(t, err)
	assert.True(t, allowed)
}
