package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout/internal/pricing"
)

func main() {
	client := pricing.FixedDelayClient{
		Delay: 50 * time.Millisecond,
		Price: pricing.Price{SKU: "sku-42", Cents: 1299},
	}
	service := pricing.NewService(client, time.Nanosecond)

	_, err := service.Lookup(context.Background(), "sku-42")
	fmt.Printf("lookup timed out: %v\n", errors.Is(err, context.DeadlineExceeded))
	fmt.Println("run make lint to expose the blocked worker send without waiting on wall-clock time")
}
