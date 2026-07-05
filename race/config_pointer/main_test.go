package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestConfigRefreshWhileReading(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(oldProcs)

	cache := NewConfigCache(&Config{APIHost: "api-0.internal"})
	configs := make([]*Config, 1000)
	for i := range configs {
		configs[i] = &Config{APIHost: fmt.Sprintf("api-%d.internal", i)}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, config := range configs {
			cache.Refresh(config)
			runtime.Gosched()
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			_ = cache.APIHost()
			runtime.Gosched()
		}
	}
}
