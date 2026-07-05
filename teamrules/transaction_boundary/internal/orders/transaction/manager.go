package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) Manager {
	return Manager{db: db}
}

func (m Manager) Within(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
