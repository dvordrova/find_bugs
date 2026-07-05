package main

import (
	"context"
	"fmt"

	"github.com/dvordrova/find_bugs/teamrules/no_panic_in_service_path/internal/payments/service"
)

func main() {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Printf("service panicked: %v\n", recovered)
		}
	}()

	payments := service.NewPaymentService(100)
	_ = payments.Capture(context.Background(), service.Payment{
		ID:     "pay_42",
		Amount: 250,
	})
}
