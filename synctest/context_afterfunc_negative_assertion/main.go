package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type AuditSink struct {
	mu      sync.Mutex
	records []string
	seen    chan struct{}
}

func NewAuditSink() *AuditSink {
	return &AuditSink{
		seen: make(chan struct{}, 1),
	}
}

func (s *AuditSink) Record(accountID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records = append(s.records, accountID)
	select {
	case s.seen <- struct{}{}:
	default:
	}
}

func (s *AuditSink) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.records)
}

func (s *AuditSink) Wait(ctx context.Context) bool {
	select {
	case <-s.seen:
		return true
	case <-ctx.Done():
		return false
	}
}

func NotifyWhenCanceled(ctx context.Context, sink *AuditSink, accountID string) func() bool {
	go sink.Record(accountID) // BUG: audit is written during registration, before cancellation.

	return context.AfterFunc(ctx, func() {
		sink.Record(accountID)
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sink := NewAuditSink()

	stop := NotifyWhenCanceled(ctx, sink, "tenant-42")
	defer stop()

	waitCtx, stopWait := context.WithTimeout(context.Background(), time.Second)
	_ = sink.Wait(waitCtx)
	stopWait()
	fmt.Printf("before cancel: %d audit record(s)\n", sink.Count())

	cancel()
	waitCtx, stopWait = context.WithTimeout(context.Background(), time.Second)
	_ = sink.Wait(waitCtx)
	stopWait()
	fmt.Printf("after cancel: %d audit record(s)\n", sink.Count())
}
