package limiter

import (
	"context"
	"time"
)

// RateLimiter é a interface que define os métodos que nosso rate limiter deve implementar.
// type RateLimiter interface {
// 	Allow(ctx context.Context, key string) (bool, error)
// 	SetLimit(limit int)
// }

type RateLimiter struct {
	store      RateLimiterStore
	rate       int
	expiration time.Duration
}

// Config contém as configurações do rate limiter.
type Config struct {
	RateLimit     int
	BlockDuration time.Duration
}

func NewRateLimiter(store RateLimiterStore, rate int, expiration time.Duration) *RateLimiter {
	return &RateLimiter{
		store:      store,
		rate:       rate,
		expiration: expiration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string, limit int) bool {
	count, err := rl.store.Get(key)
	if err != nil {
		return false
	}

	if count >= limit {
		return false
	}

	_, err = rl.store.Increment(key)
	if err != nil {
		return false
	}

	err = rl.store.Set(key, count+1, int(rl.expiration.Seconds()))
	return err == nil
}
