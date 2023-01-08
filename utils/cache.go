package utils

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func NewCache(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	// Create new cache with specified default expiration and cleanup interval
	c := cache.New(defaultExpiration, cleanupInterval)
	return c
}

func SetCache(c *cache.Cache, key string, value interface{}, expiration time.Duration) {
	// Add key-value pair to cache with specified expiration
	c.Set(key, value, expiration)
}

func GetCache(c *cache.Cache, key string) (interface{}, bool) {
	// Retrieve value from cache for specified key
	value, found := c.Get(key)
	return value, found
}
