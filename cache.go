package stormpath

import (
	"encoding/json"
	"fmt"
	"errors"
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
// Custom Implmentation of Cache that works with AppEngine
//////////////////////

// A wrapper for a cache.Cache that provides Cachable methods to work with the StormPath client.
type CacheableCacheWrapper struct {
	// The wrapped cache.Cache
	Cache *cache.Cache
}

func (cc *CacheableCacheWrapper) Exists(key string) bool {
	_, exists := cc.Cache.Get(key);
	return exists
}

func (cc *CacheableCacheWrapper) Set(key string, data interface{}) {
	if valueAsJson, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		// Uses the default value defined in the wrappedCache.
		cc.Cache.Set(key, string(valueAsJson), 0)
	}
}

func (cc *CacheableCacheWrapper) Get(key string, result interface{}) error {
	_result, exists := cc.Cache.Get(key);
	if exists {
		//log.Printf("Value found for Key %v!", key)
		json.Unmarshal([]byte(_result.(string)), result)
		//log.Printf("AFTER UNMARSHAL: %#v", result)
		return nil
	} else {
		return errors.New(fmt.Sprintf("No entry found for Key '%v'", key));
	}
}

func (cc *CacheableCacheWrapper) Del(key string) {
	cc.Cache.Delete(key)
}
