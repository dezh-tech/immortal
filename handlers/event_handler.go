package handlers

import (
	"context"
	"log"
	"time"

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
				eTags = append(eTags, t[1])
				continue
			}

			if t[0] == "p" {
				pTags = append(pTags, t[1])
			}
		}

		textNote := dbmodels.TextNote{
			ID:                  e.ID,
			Content:             null.StringFrom(e.Content),
			UsersMetadatapubKey: null.StringFrom(e.PublicKey),
			EventCreatedAt:      time.Unix(e.CreatedAt, 0),
			Event:               e.String(),
			ETags:                   eTags,
			PTags:                   pTags,
		}
		err := textNote.InsertG(context.Background(), boil.Infer())
		if err != nil {
			log.Print(err)
		}
	} else if e.Kind == types.KindReaction {
		eTags := make([]string, 0)
		pTags := make([]string, 0)
		aTags := make([]string, 0)
		kTags := make([]string, 0)
		rTags := make([]string, 0)

		for _, t := range e.Tags {
			if len(t) < 2 {
				continue
			}

			if t[0] == "e" {
				eTags = append(eTags, t[1])
				continue
			}

			if t[0] == "p" {
				pTags = append(pTags, t[1])
				continue
			}

			if t[0] == "a" {
				aTags = append(aTags, t[1])
				continue
			}

			if t[0] == "k" {
				kTags = append(kTags, t[1])
				continue
			}

			if t[0] == "r" {
				rTags = append(rTags, t[1])
			}
		}

		reaction := dbmodels.Reaction{
			ID:                  e.ID,
			UsersMetadatapubKey: null.StringFrom(e.PublicKey),
			TextNotesid:         null.StringFrom(eTags[len(eTags)-1]),
			EventCreatedAt:      time.Unix(e.CreatedAt, 0),
			Event:               e.String(),
			Content:             null.StringFrom(e.Content),
			ETags:                   eTags,
			PTags:                   pTags,
			ATags:                   aTags,
			KTags:                   kTags,
			RTags:                   rTags,
		}
		err := reaction.InsertG(context.Background(), boil.Infer())
		if err != nil {
			log.Print(err)
		}

	}
}
