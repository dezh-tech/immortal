package handler

import (
	"context"
	"errors"
	"time"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types/event"
)

func (h *Handler) HandleEvent(e *event.Event) error {
	collName, ok := database.KindToCollectionName[e.Kind]
	if !ok {
		return errors.New("invalid kind")
	}
	coll := h.DB.Client.Database(h.DB.DBName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	_, err := coll.InsertOne(ctx, e)
	if err != nil {
		cancel()
		// TODO ::: refactor errors
		return err
	}
	cancel()

	return nil
}
