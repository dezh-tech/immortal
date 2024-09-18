package handlers

import (
	"context"
	"errors"
	"time"

	dbmodels "github.com/dezh-tech/immortal/database/models"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (eh *EventHandler) handleTextNote(e *event.Event) error {
	eTags := make([]string, 0)
	pTags := make([]string, 0)

	for _, t := range e.Tags {
		if len(t) < 2 {
			continue
		}

		switch t[0] {
		case "e":
			eTags = append(eTags, t[1])
		case "p":
			pTags = append(pTags, t[1])
		}
	}

	textNote := dbmodels.TextNote{
		ID:                  e.ID,
		UsersMetadatapubKey: null.StringFrom(e.PublicKey),
		EventCreatedAt:      time.Unix(e.CreatedAt, 0),
		Event:               e.String(),
		ETags:               eTags,
		PTags:               pTags,
	}
	err := textNote.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (eh *EventHandler) handleReaction(e *event.Event) error {
	eTags := make([]string, 0)
	pTags := make([]string, 0)
	aTags := make([]string, 0)
	kTags := make([]string, 0)
	rTags := make([]string, 0)

	for _, t := range e.Tags {
		if len(t) < 2 {
			continue
		}

		switch t[0] {
		case "e":
			eTags = append(eTags, t[1])
		case "p":
			pTags = append(pTags, t[1])
		case "a":
			aTags = append(aTags, t[1])
		case "k":
			kTags = append(kTags, t[1])
		case "r":
			rTags = append(rTags, t[1])
		}
	}

	var textNotesID null.String
	if len(eTags) > 0 {
		textNotesID = null.StringFrom(eTags[len(eTags)-1])
	} else {
		return errors.New("invalid: reaction")
	}

	reaction := dbmodels.Reaction{
		ID:                  e.ID,
		UsersMetadatapubKey: null.StringFrom(e.PublicKey),
		TextNotesid:         textNotesID,
		EventCreatedAt:      time.Unix(e.CreatedAt, 0),
		Event:               e.String(),
		Content:             null.StringFrom(e.Content),
		ETags:               eTags,
		PTags:               pTags,
		ATags:               aTags,
		KTags:               kTags,
		RTags:               rTags,
	}
	err := reaction.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (eh *EventHandler) handleFollowList(e *event.Event) error {
	usermetadata := dbmodels.UsersMetadatum{
		PubKey:          e.PublicKey,
		FollowListEvent: null.StringFrom(e.String()),
	}

	// TODO ::: update the follow_list relations
	err := usermetadata.UpsertG(context.Background(), true, []string{"pub_key"},
		boil.Whitelist("follow_list_event"), boil.Infer())
	if err != nil {
		return err
	}

	return nil
}
