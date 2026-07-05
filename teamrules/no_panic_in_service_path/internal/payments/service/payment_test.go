package service

import (
	"context"
	"testing"
)

func TestCaptureApprovedPayment(t *testing.T) {
	payments := NewPaymentService(100)

	err := payments.Capture(context.Background(), Payment{
		ID:     "pay_ok",
		Amount: 75,
	})
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
}
