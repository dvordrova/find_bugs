package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Receipt struct {
	Address    string
	DeliveryID string
}

type EmailGateway struct {
	Latency time.Duration
}

func SendWelcomeEmail(ctx context.Context, gateway EmailGateway, address string) error {
	receipts := gateway.Deliver(address)

	select {
	case receipt := <-receipts:
		fmt.Printf("queued welcome email %s for %s\n", receipt.DeliveryID, receipt.Address)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g EmailGateway) Deliver(address string) <-chan Receipt {
	receipts := make(chan Receipt)

	go func() {
		time.Sleep(g.Latency)
		receipts <- Receipt{
			Address:    address,
			DeliveryID: "welcome-001",
		}
	}()

	return receipts
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	gateway := EmailGateway{Latency: 50 * time.Millisecond}
	if err := SendWelcomeEmail(ctx, gateway, "alice@example.com"); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("welcome email timed out")
			return
		}
		fmt.Printf("welcome email failed: %v\n", err)
	}
}
