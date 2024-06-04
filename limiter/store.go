package limiter

type RateLimiterStore interface {
	Get(key string) (int, error)
	Set(key string, value int, expiration int) error
	Increment(key string) (int, error)
}
