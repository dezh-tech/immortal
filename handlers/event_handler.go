package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types/event"
)

type EventHandler struct {
	DB *database.Database
}

func NewEventHandler(db *database.Database) *EventHandler {
	return &EventHandler{
		DB: db,
	}
}

func (eh *EventHandler) Handle(e *event.Event) error {
	collName, ok := database.KindToCollectionName[e.Kind]
	if !ok {
		return errors.New("invalid kind")
	}
	coll := eh.DB.Client.Database("immortal_dev").Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	_, err := coll.InsertOne(ctx, e)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}
