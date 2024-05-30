package middleware

import (
	"net/http"

	"go-rate-limiter/limiter"

	"github.com/gin-gonic/gin"
)

// RateLimiterMiddleware cria um middleware de rate limiting.
func RateLimiterMiddleware(rl limiter.RateLimiter, tokenLimit int, ipLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		clientIP := c.ClientIP()
		apiKey := c.GetHeader("API_KEY")

		if apiKey != "" {
			rl.SetLimit(tokenLimit)
			key := "token:" + apiKey
			allowed, err := rl.Allow(ctx, key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				c.Abort()
				return
			}
			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
				c.Abort()
				return
			}
		} else {
			rl.SetLimit(ipLimit)
			key := "ip:" + clientIP
			allowed, err := rl.Allow(ctx, key)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				c.Abort()
				return
			}
			if !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
