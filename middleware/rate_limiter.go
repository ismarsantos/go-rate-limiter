package middleware

import (
	"net/http"

	"go-rate-limiter/limiter"

	"github.com/gin-gonic/gin"
)

// RateLimiterMiddleware cria um middleware de rate limiting.
func RateLimiterMiddleware(rl *limiter.RateLimiter, tokenLimit int, ipLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIP := c.ClientIP()
		apiKey := c.GetHeader("API_KEY")

		var key string
		var limit int

		if apiKey != "" {
			key = "token:" + apiKey
			limit = tokenLimit
		} else {
			key = "ip:" + clientIP
			limit = ipLimit
		}

		allowed := rl.Allow(ctx, key, limit)
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			c.Abort()
			return
		}

		c.Next()
	}
}
