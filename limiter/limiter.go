package limiter

import (
	"context"
	"fmt"
	"time"
)

// RateLimiter é a interface que define os métodos que nosso rate limiter deve implementar.
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	SetLimit(limit int)
}

// Config contém as configurações do rate limiter.
type Config struct {
	RateLimit     int
	BlockDuration time.Duration
}

// NewRateLimiter cria uma nova instância do RateLimiter com base na estratégia fornecida.
func NewRateLimiter(strategy string, config Config) (RateLimiter, error) {
	switch strategy {
	case "redis":
		return NewRedisRateLimiter(config), nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strategy)
	}
}
