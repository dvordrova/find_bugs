package main

import (
	"context"
	"testing"
	"time"
)

func TestSessionCacheWarmStoresInitialValue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := NewSessionCache()
	cache.Warm(ctx, "tenant-acme", time.Hour)

	if _, ok := cache.LastRefresh("tenant-acme"); !ok {
		t.Fatal("expected session cache to store initial refresh time")
	}

	cancel()
	time.Sleep(10 * time.Millisecond)
}
