package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRateLimiter é a implementação do RateLimiter usando o Redis.
type RedisStore struct {
	client *redis.Client
	config Config
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

// Allow verifica se a requisição é permitida com base no IP ou token.
func (rl *RedisStore) Allow(ctx context.Context, key string) (bool, error) {
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
func (r *RedisStore) Set(key string, value int, expiration int) error {
	return r.client.Set(context.Background(), key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisStore) Get(key string) (int, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *RedisStore) Increment(key string) (int, error) {
	val, err := r.client.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
