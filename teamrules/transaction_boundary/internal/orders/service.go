package orders

import (
	"context"
	"database/sql"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return Service{db: db}
}

func (s Service) Capture(ctx context.Context, orderID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin capture transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, "UPDATE orders SET captured = TRUE WHERE id = ?", orderID); err != nil {
		return fmt.Errorf("mark order captured: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO outbox(order_id, event_type) VALUES (?, ?)", orderID, "order.captured"); err != nil {
		return fmt.Errorf("enqueue capture event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit capture transaction: %w", err)
	}
	return nil
}
