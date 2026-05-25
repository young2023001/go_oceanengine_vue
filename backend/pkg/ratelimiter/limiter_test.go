package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
}

func TestTenantRateLimiter_Allow(t *testing.T) {
	rdb := setupTestRedis()
	ctx := context.Background()
	rdb.FlushDB(ctx)
	defer rdb.FlushDB(ctx)

	limiter := NewTenantRateLimiter(rdb, 10)

	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, allowed, "request %d should be allowed", i+1)
	}

	allowed, err := limiter.Allow(ctx, 1)
	assert.NoError(t, err)
	assert.False(t, allowed, "request 11 should be denied")

	time.Sleep(1100 * time.Millisecond)
	allowed, err = limiter.Allow(ctx, 1)
	assert.NoError(t, err)
	assert.True(t, allowed, "request after 1s should be allowed")
}

func TestTenantRateLimiter_IsolatesTenants(t *testing.T) {
	rdb := setupTestRedis()
	ctx := context.Background()
	rdb.FlushDB(ctx)
	defer rdb.FlushDB(ctx)

	limiter := NewTenantRateLimiter(rdb, 2)

	limiter.Allow(ctx, 1)
	limiter.Allow(ctx, 1)
	allowed, _ := limiter.Allow(ctx, 1)
	assert.False(t, allowed, "tenant 1 should be denied")

	allowed, _ = limiter.Allow(ctx, 2)
	assert.True(t, allowed, "tenant 2 should not be affected")
}
