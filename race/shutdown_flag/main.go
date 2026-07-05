package main

import (
	"fmt"
	"time"
)

type Worker struct {
	stopping bool
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Run(pollEvery time.Duration) {
	ticker := time.NewTicker(pollEvery)
	defer ticker.Stop()

	for !w.stopping {
		<-ticker.C
	}
}

func (w *Worker) Stop() {
	w.stopping = true
}

func main() {
	worker := NewWorker()
	done := make(chan struct{})

	go func() {
		defer close(done)
		worker.Run(5 * time.Millisecond)
	}()

	time.Sleep(20 * time.Millisecond)
	worker.Stop()
	<-done

	fmt.Println("worker stopped")
}
