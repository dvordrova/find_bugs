package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Invoice struct {
	ID         int64
	Customer   string
	TotalCents int64
}

type InvoiceStore struct {
	db *sql.DB
}

func (s InvoiceStore) OpenInvoices(ctx context.Context) ([]Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `
		select id, customer, total_cents
		from invoices
		where paid_at is null
	`)
	if err != nil {
		return nil, err
	}

	invoices := make([]Invoice, 0, 16)
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Customer, &invoice.TotalCents); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

func main() {
	fmt.Println("load open invoices")
}
