package repository

import (
	"context"

	"github.com/dezh-tech/immortal/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) DeleteByID(id string, kind types.Kind) error {
	coll := h.db.Client.Database(h.db.DBName).Collection(getCollectionName(kind))

	ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)
	defer cancel()

	filter := bson.D{
		{Key: "id", Value: id},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "id", Value: id},
		}},
		{Key: "$unset", Value: bson.D{
			{Key: "pubkey"},
			{Key: "created_at"},
			{Key: "kind"},
			{Key: "tags"},
			{Key: "content"},
			{Key: "sig"},
		}},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
