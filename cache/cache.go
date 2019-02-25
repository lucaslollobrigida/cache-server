package cache

import (
	"fmt"
	"time"
)

type Cache struct {
	Map map[string]*Registry
}

type Registry struct {
	RegTime time.Time `json:"regTime"`
	Value   string    `json:"value"`
}

func (c *Cache) Insert(key string, value string) {
	reg := Registry{
		time.Now(),
		value,
	}

	fmt.Println(key, reg)
	c.Map[key] = &reg
}

func (c *Cache) Remove(key string) {
	delete(c.Map, key)
}

func (c *Cache) Get(key string) *Registry {
	return c.Map[key]
}
