package handlers

import (
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types"
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
	switch e.Kind { //nolint
	case types.KindTextNote:
		return eh.handleTextNote(e)
	case types.KindReaction:
		return eh.handleReaction(e)
	case types.KindFollowList:
		return eh.handleFollowList(e)
	default:
		return nil
	}
}
