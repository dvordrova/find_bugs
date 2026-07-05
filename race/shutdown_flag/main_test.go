package main

import (
	"runtime"
	"testing"
	"time"
)

func TestStopWhileWorkerChecksFlag(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(oldProcs)

	worker := NewWorker()
	done := make(chan struct{})

	go func() {
		defer close(done)
		worker.Run(100 * time.Microsecond)
	}()

	time.Sleep(time.Millisecond)
	worker.Stop()
	<-done
}
