package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" //nolint
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Database struct {
	db *sql.DB
}

func New(cfg Config) (*Database, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	boil.SetDB(db)

	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.MaxConnLifeTime)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxIdleTime(cfg.MaxIdleConnTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (db *Database) Stop() error {
	return db.db.Close()
}
