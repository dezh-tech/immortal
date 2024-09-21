package handlers

import (
	"context"
	"time"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReqHandler struct {
	DB *database.Database
}

type queryFilter struct {
	Tags map[string]types.Tag

	Authors []string
	IDs     []string

	Since int64
	Until int64
	Limit uint16
}

func NewReqHandler(db *database.Database) *ReqHandler {
	return &ReqHandler{
		DB: db,
	}
}

func (rh *ReqHandler) Handle(fs filter.Filters) ([]event.Event, error) {
	queryKinds := make(map[types.Kind][]queryFilter)

	for _, f := range fs {
		qf := queryFilter{
			Tags:    f.Tags,
			Authors: f.Authors,
			IDs:     f.IDs,
			Since:   f.Since,
			Until:   f.Until,
			Limit:   f.Limit,
		}

		uniqueKinds := removeDuplicateKind(f.Kinds)
		for _, k := range uniqueKinds {
			queryKinds[k] = append(queryKinds[k], qf)
		}
	}

	var finalResult []event.Event

	for kind, filters := range queryKinds {
		collection := rh.DB.Client.Database("immortal_dev").Collection(database.KindToCollectionName[kind])
		for _, f := range filters {
			query, opts, err := BuildMongoDynamicQueries(f)
			if err != nil {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			cursor, err := collection.Find(ctx, query, opts)
			if err != nil {
				return nil, err
			}

			var results []event.Event
			if err = cursor.All(context.TODO(), &results); err != nil {
				return nil, err
			}

			finalResult = append(finalResult, results...)
		}
	}

	return finalResult, nil
}

func removeDuplicateKind(intSlice []types.Kind) []types.Kind {
	allKeys := make(map[types.Kind]bool)
	list := []types.Kind{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func BuildMongoDynamicQueries(filter queryFilter) (bson.D, *options.FindOptions, error) {
	var query bson.D
	opts := options.Find()

	// Filter by IDs
	if len(filter.IDs) > 0 {
		query = append(query, bson.E{Key: "id", Value: bson.M{"$in": filter.IDs}})
	}

	// Filter by Authors
	if len(filter.Authors) > 0 {
		query = append(query, bson.E{Key: "pubkey", Value: bson.M{"$in": filter.Authors}})
	}

	// Filter by Tags
	if len(filter.Tags) > 0 {
		tagQueries := bson.A{}
		for tagKey, tagValues := range filter.Tags {
			filter := bson.M{
				"tags": bson.M{
					"$elemMatch": bson.M{
						"0": tagKey,
						"1": bson.M{"$in": tagValues},
					},
				},
			}
			tagQueries = append(tagQueries, filter)
		}
			query = append(query, bson.E{Key: "$and", Value: tagQueries})
	}

	// Filter by Since (created_at >=)
	if filter.Since > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$gte": time.Unix(filter.Since, 0)}})
	}

	// Filter by Until (created_at <=)
	if filter.Until > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$lte": time.Unix(filter.Until, 0)}})
	}

	// Add Limit to options
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}

	return query, opts, nil
}
