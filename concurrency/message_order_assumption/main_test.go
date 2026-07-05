package main

import (
	"testing"

	"github.com/dvordrova/find_bugs/concurrency/message_order_assumption/internal/projection"
)

func TestOrdinaryProcessesBrokerOrder(t *testing.T) {
	projector := projection.New()
	projector.Apply(projection.AccountOpened{AccountID: "acct-42"})
	projector.Apply(projection.CreditReserved{AccountID: "acct-42", Cents: 500})

	if got := projector.ReservedCents("acct-42"); got != 500 {
		t.Fatalf("reserved cents = %d, want 500", got)
	}
}

func TestMessageOrderPermutations(t *testing.T) {
	tests := []struct {
		name   string
		events []any
	}{
		{
			name: "account_before_credit",
			events: []any{
				projection.AccountOpened{AccountID: "acct-42"},
				projection.CreditReserved{AccountID: "acct-42", Cents: 500},
			},
		},
		{
			name: "credit_before_account",
			events: []any{
				projection.CreditReserved{AccountID: "acct-42", Cents: 500},
				projection.AccountOpened{AccountID: "acct-42"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projector := projection.New()
			for _, event := range tt.events {
				projector.Apply(event)
			}

			if got := projector.ReservedCents("acct-42"); got != 500 {
				t.Fatalf("reserved cents after %s = %d, want 500", tt.name, got)
			}
		})
	}
}
