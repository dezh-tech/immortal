package database

import (
	"database/sql"

	_ "github.com/lib/pq" // no-lint
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Database struct {
	db *sql.DB
}

func New(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	boil.SetDB(db)
	// TODO ::: config the connection pool
	return &Database{
		db: db,
	}, nil
}
