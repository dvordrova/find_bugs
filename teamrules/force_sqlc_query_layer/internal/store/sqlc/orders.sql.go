package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Order struct {
	ID     string
	Status string
}

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) ListPendingOrders(ctx context.Context, limit int) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, `
		SELECT id, status
		FROM orders
		WHERE status = ?
		ORDER BY created_at
		LIMIT ?
	`, "pending", limit)
	if err != nil {
		return nil, fmt.Errorf("query pending orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Status); err != nil {
			return nil, fmt.Errorf("scan pending order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending orders: %w", err)
	}
	return orders, nil
}
