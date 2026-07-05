package main

import (
	"fmt"
	"sync"
)

type AccountLedger struct {
	mu       sync.Mutex
	balances map[string]int
}

func NewAccountLedger() AccountLedger {
	return AccountLedger{
		balances: map[string]int{
			"alice": 120,
		},
	}
}

func (l AccountLedger) Balance(accountID string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.balances[accountID]
}

func main() {
	ledger := NewAccountLedger()

	fmt.Printf("alice balance: %d\n", ledger.Balance("alice"))
}
