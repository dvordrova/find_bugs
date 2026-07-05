package customer

import (
	"context"
	"database/sql"
)

type Customer struct {
	ID    int64
	Email string
}

type ExportService struct {
	db *sql.DB
}

func NewExportService(db *sql.DB) ExportService {
	return ExportService{db: db}
}

func (s ExportService) ActiveCustomers(ctx context.Context) ([]Customer, error) {
	rows, err := s.db.QueryContext(ctx, `
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
