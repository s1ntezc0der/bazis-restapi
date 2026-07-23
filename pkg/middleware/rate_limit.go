package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	ctx := context.Background()

	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		rl.client.Expire(ctx, key, rl.window)
	}

	return count <= int64(rl.limit)
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(UserIDKey)
			var key string

			if userID != nil {
				key = "rate_limit:user:" + strconv.Itoa(userID.(int))
			} else {
				key = "rate_limit:ip:" + r.RemoteAddr
			}

			if !rl.Allow(key) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GinRateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		var key string
		if userID != nil {
			key = "rate_limit:user:" + strconv.FormatInt(userID.(int64), 10)
		} else {
			key = "rate_limit:ip:" + c.ClientIP()
		}

		if !rl.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
