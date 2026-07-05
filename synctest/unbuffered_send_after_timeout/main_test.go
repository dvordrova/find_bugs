package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/dvordrova/find_bugs/synctest/unbuffered_send_after_timeout/internal/pricing"
)

func TestOrdinaryTestOnlyChecksTimeout(t *testing.T) {
	client := pricing.FixedDelayClient{
		Delay: time.Hour,
		Price: pricing.Price{SKU: "sku-42", Cents: 1299},
	}
	service := pricing.NewService(client, time.Nanosecond)

	_, err := service.Lookup(context.Background(), "sku-42")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Lookup() error = %v, want DeadlineExceeded", err)
	}
}

func TestSynctestFindsBlockedSendAfterTimeout(t *testing.T) {
	recovered := runBlockedSendCheck(t)
	if recovered == nil {
		t.Fatal("synctest check returned without detecting a blocked worker send")
	}

	message := fmt.Sprint(recovered)
	if !strings.Contains(message, "blocked goroutines remain") {
		t.Fatalf("unexpected synctest panic: %s", message)
	}
	t.Fatalf("synctest detected blocked worker send: %s", message)
}

func runBlockedSendCheck(t *testing.T) (recovered any) {
	t.Helper()
	defer func() {
		recovered = recover()
	}()

	synctest.Test(t, func(t *testing.T) {
		const timeout = time.Second
		client := pricing.FixedDelayClient{
			Delay: 2 * timeout,
			Price: pricing.Price{SKU: "sku-42", Cents: 1299},
		}
		service := pricing.NewService(client, timeout)

		_, err := service.Lookup(t.Context(), "sku-42")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Lookup() error = %v, want DeadlineExceeded", err)
		}

		time.Sleep(timeout)
		synctest.Wait()
	})
	return nil
}
