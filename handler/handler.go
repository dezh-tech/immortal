package handler

import "github.com/dezh-tech/immortal/database"

type Handler struct {
	DB *database.Database
}

func New(db *database.Database) Handler {
	return Handler{
		DB: db,
	}
}
