package service

import (
	"context"
	"fmt"
)

type Payment struct {
	ID     string
	Amount int64
}

type PaymentService struct {
	limit int64
}

func NewPaymentService(limit int64) PaymentService {
	return PaymentService{limit: limit}
}

func (s PaymentService) Capture(ctx context.Context, payment Payment) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if payment.Amount > s.limit {
		panic(fmt.Sprintf("payment %s exceeds capture limit", payment.ID))
	}

	return nil
}
