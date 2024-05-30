package limiter

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// RedisRateLimiter é a implementação do RateLimiter usando o Redis.
type RedisRateLimiter struct {
	client *redis.Client
	config Config
}

// NewRedisRateLimiter cria uma nova instância do RedisRateLimiter.
func NewRedisRateLimiter(config Config) *RedisRateLimiter {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &RedisRateLimiter{
		client: client,
		config: config,
	}
}

// Allow verifica se a requisição é permitida com base no IP ou token.
func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	// Incrementa o contador e obtém o valor atual
	val, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Se for a primeira requisição, define o TTL
	if val == 1 {
		rl.client.Expire(ctx, key, rl.config.BlockDuration).Result()
	}

	// Verifica se o valor excede o limite
	if val > int64(rl.config.RateLimit) {
		return false, nil
	}

	return true, nil
}

// SetLimit ajusta o limite de requisições permitidas.
func (rl *RedisRateLimiter) SetLimit(limit int) {
	rl.config.RateLimit = limit
}
