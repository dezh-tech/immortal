package repository

import (
	"context"

	"github.com/dezh-tech/immortal/pkg/logger"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleReq(f *filter.Filter) ([]event.Event, error) {
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

	var pipeline mongo.Pipeline

	for kind, filter := range queryKinds {
		collectionName, isMultiKindColl := getCollectionName(kind)

		query := filterToMongoQuery(filter, isMultiKindColl, kind)

		matchStage := bson.D{
			{Key: "$match", Value: query},
		}

		unionStage := bson.D{
			{Key: "$unionWith", Value: bson.D{
				{Key: "coll", Value: collectionName},
				{Key: "pipeline", Value: mongo.Pipeline{
					matchStage,
				}},
			}},
		}

		pipeline = append(pipeline, unionStage)
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "created_at", Value: -1},
			{Key: "id", Value: 1},
		}},
	}

	pipeline = append(pipeline, sortStage)

	finalLimit := h.config.DefaultQueryLimit
	if f.Limit > 0 && f.Limit < h.config.MaxQueryLimit {
		finalLimit = uint32(f.Limit)
	}

	limitStage := bson.D{
		{Key: "$limit", Value: finalLimit},
	}

	pipeline = append(pipeline, limitStage)

	ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)
	defer cancel()

	collection := h.db.Client.Database(h.db.DBName).Collection("empty")
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		_, err := h.grpc.AddLog(context.Background(),
			"database error while adding new event", err.Error())
		if err != nil {
			logger.Error("can't send log to manager", "err", err)
		}

		return nil, err

	}
	defer cursor.Close(ctx)

	var finalResult []event.Event
	if err := cursor.All(ctx, &finalResult); err != nil {
		_, err := h.grpc.AddLog(context.Background(),
			"database error while adding new event", err.Error())
		if err != nil {
			logger.Error("can't send log to manager", "err", err)
		}

		return nil, err

	}

	return finalResult, nil
}
