package stormpath_test

import (
	. "github.com/sappenin/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickmn/go-cache"
	"time"
)

var _ = Describe("Cache", func() {
	Describe("CacheWrapper", func() {
		key := "key"

		var c *cache.Cache = cache.New(5 * time.Minute, 30 * time.Second)
		cache := &CacheableCache{Cache: c}
		AfterEach(func() {
			cache.Cache.Flush()
		})

		Describe("Exists", func() {
			It("should return false if the key doesn't exists", func() {
				r := cache.Exists(key)

				Expect(r).To(BeFalse())
			})
			It("should return true if the key does exists", func() {
				cache.Cache.Set(key, "hello", 5 * time.Minute)

				r := cache.Exists(key)

				Expect(r).To(BeTrue())
			})
		})

		Describe("Set", func() {
			It("should store a new object in the cache", func() {
				r := cache.Exists(key)

				Expect(r).To(BeFalse())

				cache.Set(key, "hello")

				r = cache.Exists(key)

				Expect(r).To(BeTrue())
			})
			It("should update an existing object in the cache", func() {
				cache.Set(key, "hello")
				cache.Set(key, "bye")

				var r string
				cache.Get(key, &r)

				//Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal("bye"))
			})
		})
		Describe("Get", func() {
			It("should load empty data if the key doesn't exists into the given interface", func() {
				var r string

				err := cache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(BeEmpty())
			})
			It("should load data from the cache into the given interface", func() {
				var r string

				cache.Set(key, "hello")
				err := cache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal("hello"))
			})
		})
		Describe("Del", func() {
			It("should delete a given key from the cache", func() {
				cache.Set(key, []byte("hello"))

				r := cache.Exists(key)
				Expect(r).To(BeTrue())

				cache.Del(key)

				r = cache.Exists(key)
				Expect(r).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Cacheable", func() {
	Describe("Collection resource", func() {
		It("should not be cacheable", func() {
			var resources = []interface{}{
				&Applications{},
				&Accounts{},
				&Groups{},
				&Directories{},
				&AccountStoreMappings{},
			}
			for _, resource := range resources {
				c, ok := resource.(Cacheable)

				Expect(ok).To(BeTrue())
				Expect(c.IsCacheable()).To(BeFalse())
			}
		})
	})

	Describe("Single resource", func() {
		It("should be cacheable", func() {
			var resources = []interface{}{
				&Application{},
				&Account{},
				&Group{},
				&Directory{},
				&AccountStoreMapping{},
				&Tenant{},
			}
			for _, resource := range resources {
				c, ok := resource.(Cacheable)

				Expect(ok).To(BeTrue())
				Expect(c.IsCacheable()).To(BeTrue())
			}
		})
	})
})