package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewReportCache(rdb *redis.Client) *ReportCache {
	return &ReportCache{rdb: rdb, ttl: 5 * time.Minute}
}

func (c *ReportCache) Get(ctx context.Context, key string, dest interface{}) bool {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false
	}
	return true
}

func (c *ReportCache) Set(ctx context.Context, key string, value interface{}) {
	data, _ := json.Marshal(value)
	c.rdb.Set(ctx, key, string(data), c.ttl)
}

func (c *ReportCache) Invalidate(ctx context.Context, tenantID uint64) {
	pattern := fmt.Sprintf("report:%d:*", tenantID)
	iter := c.rdb.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		c.rdb.Del(ctx, keys...)
	}
}

func OverviewKey(tenantID uint64, date string) string {
	return fmt.Sprintf("report:%d:overview:%s", tenantID, date)
}

func TrendKey(tenantID uint64, start, end string) string {
	return fmt.Sprintf("report:%d:trend:%s_%s", tenantID, start, end)
}
