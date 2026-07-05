package domain

import (
	"database/sql"
	"net/http"
)

type Invoice struct {
	ID           string
	CustomerNote sql.NullString
	Headers      http.Header
}

func (i Invoice) Note() string {
	if !i.CustomerNote.Valid {
		return ""
	}
	return i.CustomerNote.String
}
