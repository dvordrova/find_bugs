package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return InvoiceRepository{db: db}
}

func (r InvoiceRepository) MarkPaid(ctx context.Context, invoiceID string) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE invoices SET paid = TRUE WHERE id = ?", invoiceID); err != nil {
		return fmt.Errorf("mark invoice paid: %w", err)
	}
	return nil
}
