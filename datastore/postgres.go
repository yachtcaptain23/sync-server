package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgres is a Datastore wrapper around a postgres database
type Postgres struct {
	*sqlx.DB
}

func NewPostgres() (*Postgres, error) {
	db, err := sqlx.Open("postgres", "user=yrliou dbname=sync_server sslmode=disable")
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(80)

	pg := &Postgres{db}
	return pg, nil
}
