package repository

import (
	"context"
	"strconv"

	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) DeleteByID(id string, kind types.Kind) error {
	deleteFilter := bson.D{
		{Key: "id", Value: id},
	}

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: "pubkey"},
			{Key: "created_at"},
			{Key: "kind"},
			{Key: "tags"},
			{Key: "content"},
			{Key: "sig"},
		}},
	}

	collName, _ := getCollectionName(kind)
	coll := h.db.Client.Database(h.db.DBName).Collection(collName)

	ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)
	defer cancel()

	_, err := coll.UpdateOne(ctx, deleteFilter, update)
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

func (h *Handler) NIP09Deletion(e *event.Event) error {
	kinds := e.Tags.GetValues("k")
	eventIDs := e.Tags.GetValues("e")

	queryKinds := []types.Kind{}

	for _, k := range kinds {
		k, err := strconv.ParseInt(k, 10, 16)
		if err != nil {
			continue
		}

		// todo:: update gosec linter and remove //nolint comment.
		queryKinds = append(queryKinds, types.Kind(k)) //nolint
	}

	deleteFilter := bson.D{
		{Key: "pubkey", Value: e.PublicKey},
	}

	deleteFilter = append(deleteFilter, bson.E{Key: "id", Value: bson.M{"$in": eventIDs}})

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: "pubkey"},
			{Key: "created_at"},
			{Key: "kind"},
			{Key: "tags"},
			{Key: "content"},
			{Key: "sig"},
		}},
	}

	for _, kind := range queryKinds {
		ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)

		collectionName, _ := getCollectionName(kind)
		coll := h.db.Client.Database(h.db.DBName).Collection(collectionName)

		_, err := coll.UpdateOne(ctx, deleteFilter, update)
		if err != nil {
			_, err := h.grpc.AddLog(context.Background(),
				"database error while removing event", err.Error())
			if err != nil {
				logger.Error("can't send log to manager", "err", err)
			}

			cancel()

			return err
		}
		cancel()
	}

	return nil
}

func (h *Handler) DeleteByFilter(f *filter.Filter) error {
	// question/todo::: is it possible to run the deletion on all collections with one database call?
	// we have an open issue on deletion execution.
	// we do the read operation using aggregation pipeline and $unionWith stage which
	// helps us ti prevent multiple database calls and it would help us to do the operation faster.
	// to do the same thing for deletion we need to filter the documents with $match, then update the
	// fields of deleted event to null (expect the `id` since its unique index to prevent overwrites) with $unset
	// then we apply them to collection using $merge.
	// although we can't use multiple $merge's on one pipeline and we must have
	// only one merge at the end of pipeline commands. also, $unionWith is restricted to be used with $merge.

	//  notes::: these details may help you to think for solutions better:
	// 1. we create a collection for each kind or each group of kinds.
	// using this model forces us to make query to all collections corresponding to provided kinds when
	// we are dealing with filters since filters contain a list of kinds
	// (which can be empty and we are then forced to query all collections)

	// 2. when we delete an event we $unset all fields expect `id`.
	// when we make a query to read from database, we ignore fields which
	// their fields are null. and when we write new events we prevent overwriting
	// events with duplicated `id`. so we can handle the deletion properly.

	// resources::: these links may help you:
	// 1. https://www.mongodb.com/docs/manual/reference/operator/aggregation/merge/#restrictions
	// 2. https://www.mongodb.com/docs/manual/reference/operator/aggregation/unionWith/#mongodb-pipeline-pipe.-unionWith
	// 3. https://www.mongodb.com/docs/manual/reference/operator/aggregation

	queryKinds := make(map[types.Kind]*filter.Filter)

	if len(f.Kinds) != 0 {
		uniqueKinds := removeDuplicateKinds(f.Kinds)
		for _, k := range uniqueKinds {
			queryKinds[k] = f
		}
	} else {
		for k := range types.KindToName {
			queryKinds[k] = f
		}
	}

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: "pubkey"},
			{Key: "created_at"},
			{Key: "kind"},
			{Key: "tags"},
			{Key: "content"},
			{Key: "sig"},
		}},
	}

	for kind, deleteFilter := range queryKinds {
		collectionName, isMultiKindColl := getCollectionName(kind)

		query := filterToMongoQuery(deleteFilter, isMultiKindColl, kind)

		ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)

		_, err := h.db.Client.Database(h.db.DBName).Collection(collectionName).UpdateMany(ctx, query, update)
		if err != nil {
			_, err := h.grpc.AddLog(ctx,
				"database error while deleting events", err.Error())
			if err != nil {
				logger.Error("can't send log to manager", "err", err)
			}

			cancel()

			return err
		}
		cancel()
	}

	return nil
}
