package main

import "fmt"

type Config struct {
	APIHost string
}

type ConfigCache struct {
	current *Config
}

func NewConfigCache(initial *Config) *ConfigCache {
	return &ConfigCache{current: initial}
}

func (c *ConfigCache) Refresh(next *Config) {
	c.current = next
}

func (c *ConfigCache) APIHost() string {
	current := c.current
	if current == nil {
		return ""
	}
	return current.APIHost
}

func main() {
	cache := NewConfigCache(&Config{APIHost: "api.internal"})

	fmt.Printf("api host: %s\n", cache.APIHost())
}
