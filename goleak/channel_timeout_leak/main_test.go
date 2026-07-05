package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSendWelcomeEmailReturnsDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	gateway := EmailGateway{Latency: 20 * time.Millisecond}
	err := SendWelcomeEmail(ctx, gateway, "alice@example.com")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}

	time.Sleep(30 * time.Millisecond)
}
