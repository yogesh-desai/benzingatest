package utils

import (
	cache "github.com/hashicorp/golang-lru"
)

var Cache *cache.Cache

// initCache initializes cache
func initCache(size int) (*cache.Cache, error) {

	return cache.New(size) // Get new cache as per size provided in the env variables.

}
