package handler

import (
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types"
)

var KindToCollectionName = map[types.Kind]string{
	types.KindTextNote:               "text_notes",
	types.KindReaction:               "reactions",
	types.KindProfileMetadata:        "profile_metadatas",
	types.KindFollowList:             "follow_lists",
	types.KindEncryptedDirectMessage: "encrypted_dms_nip4",
	types.KindRepost:                 "reposts",
	types.KindProfileBadges:          "badge_awards",
}

type Handler struct {
	DB *database.Database
}

func New(db *database.Database) Handler {
	return Handler{
		DB: db,
	}
}
