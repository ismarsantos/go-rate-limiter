package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"go-rate-limiter/limiter"
	"go-rate-limiter/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rateLimitIP := os.Getenv("RATE_LIMIT_IP")
	rateLimitToken := os.Getenv("RATE_LIMIT_TOKEN")
	blockDuration := os.Getenv("BLOCK_DURATION")
	webPORT := os.Getenv("PORT")

	ipLimit, err := strconv.Atoi(rateLimitIP)
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT_IP value: %v", err)
	}

	tokenLimit, err := strconv.Atoi(rateLimitToken)
	if err != nil {
		log.Fatalf("Invalid RATE_LIMIT_TOKEN value: %v", err)
	}

	blockDur, err := strconv.Atoi(blockDuration)
	if err != nil {
		log.Fatalf("Invalid BLOCK_DURATION value: %v", err)
	}

	config := limiter.Config{
		BlockDuration: time.Duration(blockDur) * time.Second,
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // Use DB 0 for tests
	})

	store := limiter.NewRedisStore(client)

	rateLimiter := limiter.NewRateLimiter(store, config.RateLimit, config.BlockDuration)
	if err != nil {
		log.Fatalf("Error creating rate limiter: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.RateLimiterMiddleware(rateLimiter, tokenLimit, ipLimit))

	// Defina as rotas do seu servidor aqui
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Go Rate Limiter Success"})
	})

	r.Run(":" + webPORT)
}
