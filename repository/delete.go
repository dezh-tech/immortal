package repository

import (
	"context"

	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) DeleteByID(id string, kind types.Kind) error {
	collName, _ := getCollectionName(kind)
	coll := h.db.Client.Database(h.db.DBName).Collection(collName)

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
		_, err := h.grpc.AddLog(context.Background(),
			"database error while removing event", err.Error())
		if err != nil {
			logger.Error("can't send log to manager", "err", err)
		}

		return err
	}

	return nil
}
