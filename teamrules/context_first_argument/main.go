package main

import (
	"context"
	"fmt"

	"github.com/dvordrova/find_bugs/teamrules/context_first_argument/internal/invoices/service"
)

func main() {
	invoices := service.NewInvoiceService()
	_ = invoices.RebuildInvoice("inv-42", context.Background())
	fmt.Println("invoice rebuild requested")
	fmt.Println("run make lint to see the revive report")
}
