package repository

import (
	"context"
	"database/sql"
)

type Customer struct {
	ID    int64
	Email string
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return Repository{db: db}
}

func (r Repository) ActiveCustomers(ctx context.Context) ([]Customer, error) {
	rows, err := r.db.QueryContext(ctx, `
		select id, email
		from customers
		where active = true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := make([]Customer, 0, 128)
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Email); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}
