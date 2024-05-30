package limiter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // Use DB 0 for tests
	})

	// Clear all keys in the test database
	client.FlushDB(context.Background())

	return client
}

func TestRedisRateLimiter_Allow(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	config := Config{
		RateLimit:     5,
		BlockDuration: 10 * time.Second,
	}

	limiter := &RedisRateLimiter{
		client: client,
		config: config,
	}

	ctx := context.Background()
	key := "test-key"

	// Test within rate limit
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// Test exceeding rate limit
	allowed, err := limiter.Allow(ctx, key)
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestRedisRateLimiter_SetLimit(t *testing.T) {
	client := setupTestRedis()
	defer client.Close()

	config := Config{
		RateLimit:     5,
		BlockDuration: 10 * time.Second,
	}

	limiter := &RedisRateLimiter{
		client: client,
		config: config,
	}

	ctx := context.Background()
	key := "test-key-set-limit"

	// Test within initial rate limit
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// Change rate limit
	limiter.SetLimit(10)

	// Test within new rate limit
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	// Test exceeding new rate limit
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		assert.NoError(t, err)
		assert.False(t, allowed)
	}
}

func TestMain(m *testing.M) {
	// Load environment variables for testing
	err := os.Setenv("REDIS_ADDR", "localhost:6379")
	if err != nil {
		panic("Error setting environment variables for tests")
	}

	os.Exit(m.Run())
}
