package stormpath

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
)

//Cacheable determines if the implementor should be cached or not
type Cacheable interface {
	IsCacheable() bool
}

//cacheResource stores a resource in the cache if the resource allows caching
func cacheResource(key string, resource interface{}, cache Cache) {
	c, ok := resource.(Cacheable)

	if ok && c.IsCacheable() {
		cache.Set(key, resource)
	}
}

//Cache is a base interface for any cache provider
type Cache interface {
	Exists(key string) bool
	Set(key string, data interface{})
	Get(key string, result interface{}) error
	Del(key string)
}

// Wrapper for

//////////////////////
// Cache implementation that works with go-cache (github.com/patrickmn/go-cache)
//////////////////////

// A wrapper for a cache.Cache that provides Cachable methods to work with the StormPath client.
type CacheableCache struct {
	// The wrapped cache.Cache
	Cache *cache.Cache
}

func (cc *CacheableCache) Exists(key string) bool {
	_, exists := cc.Cache.Get(key);
	return exists
}

func (cc *CacheableCache) Set(key string, data interface{}) {
	if valueAsJson, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		// Uses the default value defined in the wrappedCache.
		cc.Cache.Set(key, string(valueAsJson), 0)
	}
}

func (cc *CacheableCache) Get(key string, result interface{}) error {
	_result, exists := cc.Cache.Get(key);
	if !exists {
		return nil
	} else {
		return json.Unmarshal([]byte(_result.(string)), result)
	}
}

func (cc *CacheableCache) Del(key string) {
	cc.Cache.Delete(key)
}
