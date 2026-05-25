package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var luaScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local current = tonumber(redis.call('GET', key) or '0')
if current < limit then
    redis.call('INCR', key)
    redis.call('EXPIRE', key, 2)
    return 0
end
return 1
`)

type TenantRateLimiter struct {
	rdb   *redis.Client
	limit int
}

func NewTenantRateLimiter(rdb *redis.Client, qpsPerTenant int) *TenantRateLimiter {
	return &TenantRateLimiter{
		rdb:   rdb,
		limit: qpsPerTenant,
	}
}

func (l *TenantRateLimiter) Allow(ctx context.Context, tenantID uint64) (bool, error) {
	now := time.Now().Unix()
	key := fmt.Sprintf("ratelimit:%d:ocean:%d", tenantID, now)

	result, err := luaScript.Run(ctx, l.rdb, []string{key}, l.limit).Int()
	if err != nil {
		return false, err
	}
	return result == 0, nil
}

func (l *TenantRateLimiter) Wait(ctx context.Context, tenantID uint64) error {
	deadline := time.Now().Add(30 * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("rate limit wait timeout for tenant %d", tenantID)
		}
		allowed, err := l.Allow(ctx, tenantID)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}
