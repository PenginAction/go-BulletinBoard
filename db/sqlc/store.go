package db

import (
	"database/sql"
)

type Store interface {
	Querier
}

// Store provides all functions to execute DB queries
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Store provides all functions to execute SQL queries
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
