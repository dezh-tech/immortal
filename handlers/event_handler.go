package handlers

import (
	"context"
	"log"

	"github.com/dezh-tech/immortal/database"
	dbmodels "github.com/dezh-tech/immortal/database/models"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type EventHandler struct {
	DB *database.Database
}

func New(db *database.Database) *EventHandler {
	return &EventHandler{
		DB: db,
	}
}

func (eh *EventHandler) Handle(e *event.Event) {
	if e.Kind == types.KindTextNote {
		eTags := make([]string, 0)
		pTags := make([]string, 0)

		for _, t := range e.Tags {
			if len(t) < 2 {
				continue
			}

			if t[0] == "e" {
				t = t[1:]
				eTags = append(eTags, t[1])
			}

			if t[0] == "p" {
				t = t[1:]
				pTags = append(pTags, t[1])
			}
		}
		
		textNote := dbmodels.TextNote{
			ID:                  e.ID,
			Content:             null.StringFrom(e.Content),
			UsersMetadatapubKey: null.StringFrom(e.PublicKey),
			Event:               e.String(),
			E:                   eTags,
			P:                   pTags,
		}
		err := textNote.InsertG(context.Background(), boil.Infer())
		if err != nil {
			log.Fatal(err)
		}
	}
}
