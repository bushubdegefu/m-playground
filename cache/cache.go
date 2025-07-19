package cache

import (
	"strconv"
	"time"

	"github.com/bushubdegefu/m-playground/configs"
	"github.com/dgraph-io/ristretto"
)

type CacheService struct {
	cache *ristretto.Cache
}

func NewCacheService() (*CacheService, error) {
	number_of_counters, _ := strconv.ParseInt(configs.AppConfig.GetOrDefault("NUMBER_OF_COUNTERS", "250"), 10, 64)
	cache_size, _ := strconv.ParseInt(configs.AppConfig.GetOrDefault("MAX_CACHE_SIZE", "250"), 10, 64)
	buffer_items, _ := strconv.ParseInt(configs.AppConfig.GetOrDefault("BUFFER_ITEMS", "64"), 10, 64)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: number_of_counters, // Number of keys to track frequency
		MaxCost:     cache_size << 20,   // Maximum cache size (500MB)
		BufferItems: buffer_items,       // Number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return &CacheService{cache: cache}, nil
}

func (c *CacheService) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *CacheService) Set(key string, value any) bool {
	minutes, _ := strconv.Atoi(configs.AppConfig.GetOrDefault("BUFFER_ITEMS", "10"))
	ttl := time.Duration(minutes) * time.Minute
	return c.cache.SetWithTTL(key, value, 1, ttl)
}

func (c *CacheService) Delete(key string) {
	c.cache.Del(key)
}
