package main

import (
	"fmt"

	"github.com/dvordrova/find_bugs/teamrules/no_infrastructure_imports_in_domain/internal/billing/domain"
)

func main() {
	invoice := domain.Invoice{
		ID: "inv-42",
	}
	fmt.Printf("invoice note: %q\n", invoice.Note())
	fmt.Println("run make lint to see the depguard report")
}
