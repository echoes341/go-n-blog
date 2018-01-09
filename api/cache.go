package main

import (
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

// Cache from cache.go without pointer
func Cache(store persistence.CacheStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.CACHE_MIDDLEWARE_KEY, store)
		c.Next()
	}
}
