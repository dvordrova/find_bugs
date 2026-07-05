package main

import (
	"fmt"

	"github.com/dvordrova/find_bugs/concurrency/message_order_assumption/internal/projection"
)

func main() {
	inOrder := projection.New()
	inOrder.Apply(projection.AccountOpened{AccountID: "acct-42"})
	inOrder.Apply(projection.CreditReserved{AccountID: "acct-42", Cents: 500})

	reordered := projection.New()
	reordered.Apply(projection.CreditReserved{AccountID: "acct-42", Cents: 500})
	reordered.Apply(projection.AccountOpened{AccountID: "acct-42"})

	fmt.Printf("broker order reserved cents: %d\n", inOrder.ReservedCents("acct-42"))
	fmt.Printf("reordered delivery reserved cents: %d\n", reordered.ReservedCents("acct-42"))
}
