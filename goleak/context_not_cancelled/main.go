package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SessionCache struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
}

func NewSessionCache() *SessionCache {
	return &SessionCache{sessions: make(map[string]time.Time)}
}

func (c *SessionCache) Warm(ctx context.Context, tenantID string, refreshEvery time.Duration) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	c.store(tenantID)

	go func() {
		ticker := time.NewTicker(refreshEvery)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.store(tenantID)
			case <-context.Background().Done():
				return
			}
		}
	}()
}

func (c *SessionCache) LastRefresh(tenantID string) (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	refreshedAt, ok := c.sessions[tenantID]
	return refreshedAt, ok
}

func (c *SessionCache) store(tenantID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessions[tenantID] = time.Now()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := NewSessionCache()
	cache.Warm(ctx, "tenant-acme", time.Minute)

	if refreshedAt, ok := cache.LastRefresh("tenant-acme"); ok {
		fmt.Printf("session cache warmed at %s\n", refreshedAt.Format(time.RFC3339))
	}
}
