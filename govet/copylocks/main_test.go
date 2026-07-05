package main

import "testing"

func TestBalance(t *testing.T) {
	ledger := NewAccountLedger()

	if got := ledger.Balance("alice"); got != 120 {
		t.Fatalf("expected alice balance 120, got %d", got)
	}
}
