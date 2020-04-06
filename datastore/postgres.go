package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Postgres is a Datastore wrapper around a postgres database
type Postgres struct {
	*sqlx.DB
}

// NewPostgres creates a new Postgres datastore
func NewPostgres() (*Postgres, error) {
	db, err := sqlx.Open("postgres", "user=sync_server dbname=sync_server sslmode=disable")
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(80)

	pg := &Postgres{db}
	return pg, nil
}
