package main

import (
	"runtime"
	"testing"
)

func TestSnapshotWhileRecording(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(oldProcs)

	metrics := NewMetrics()
	metrics.Record("/checkout")

	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 1000 {
			metrics.Record("/checkout")
			runtime.Gosched()
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			_ = metrics.Snapshot()["/checkout"]
			runtime.Gosched()
		}
	}
}
