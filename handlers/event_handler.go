package handlers

import (
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
)

type EventHandler struct {
	DB *database.Database
}

func New(db *database.Database) *EventHandler {
	return &EventHandler{
		DB: db,
	}
}

func (eh *EventHandler) Handle(e *event.Event) error {
	if e.Kind == types.KindTextNote {
		return eh.handleTextNote(e)
	} else if e.Kind == types.KindReaction {
		return eh.handleReaction(e)
	} else if e.Kind == types.KindFollowList {
		return eh.handleFollowList(e)
	}

	return nil
}
