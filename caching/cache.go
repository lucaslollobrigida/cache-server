package caching

import (
	"sync"
	"time"
)

// Cache is a in-memory key-value store
// implemented as a hash map
type Cache struct {
	Map  map[string]*Registry
	Lock sync.Mutex
}

// Registry is a entry of the cache store.
type Registry struct {
	RegTime time.Time `json:"regTime"`
	Value   string    `json:"value"`
}

// Insert safety adds a new registry to the store
// handling concurrently operations via mutex lock.
func (c *Cache) Insert(key string, value string) {
	reg := &Registry{
		time.Now(),
		value,
	}

	c.Lock.Lock()
	c.Map[key] = reg
	c.Lock.Unlock()
}

// Remove safety removes a registry to the store
// handling concurrently operations via mutex lock.
func (c *Cache) Remove(key string) {
	c.Lock.Lock()
	delete(c.Map, key)
	c.Lock.Unlock()
}

// Get returns a registry value by key if found.
func (c *Cache) Get(key string) *Registry {
	return c.Map[key]
}

// CleanupExpired trigger a routine for cleaning expired registries from map.
func (c *Cache) CleanupExpired() {
	go func() {
		for {
			for k, v := range c.Map {
				if v.RegTime.Add(time.Minute * 1).Before(time.Now()) {
					c.Remove(k)
				}
			}
			time.Sleep(time.Second * 10)
		}
	}()
}
