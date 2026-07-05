package orders

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dvordrova/find_bugs/teamrules/force_sqlc_query_layer/internal/store/sqlc"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) ListPending(ctx context.Context, limit int) ([]sqlc.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
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

	var orders []sqlc.Order
	for rows.Next() {
		var order sqlc.Order
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
