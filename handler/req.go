package handler

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

type filterQuery struct {
	Tags map[string][]string

	Authors []string
	IDs     []string

	Since int64
	Until int64
	Limit uint16
}

func (h *Handler) HandleReq(fs filter.Filters) ([]event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.DB.QueryTimeout)
	defer cancel()

	queryKinds := make(map[types.Kind][]filterQuery)

	for _, f := range fs {
		qf := filterQuery{
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
		collection := h.DB.Client.Database(h.DB.DBName).Collection(database.KindToCollectionName[kind])
		for _, f := range filters {
			query, opts, err := FilterToQuery(&f)
			if err != nil {
				continue
			}

			cursor, err := collection.Find(ctx, query, opts)
			if err != nil {
				return nil, err
			}

			var result []event.Event
			if err = cursor.All(ctx, &result); err != nil {
				return nil, err
			}

			finalResult = append(finalResult, result...)
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

func FilterToQuery(fq *filterQuery) (bson.D, *options.FindOptions, error) {
	var query bson.D
	opts := options.Find()

	// Filter by IDs
	if len(fq.IDs) > 0 {
		query = append(query, bson.E{Key: "id", Value: bson.M{"$in": fq.IDs}})
	}

	// Filter by Authors
	if len(fq.Authors) > 0 {
		query = append(query, bson.E{Key: "pubkey", Value: bson.M{"$in": fq.Authors}})
	}

	// Filter by Tags
	if len(fq.Tags) > 0 {
		tagQueries := bson.A{}
		for tagKey, tagValues := range fq.Tags {
			qtf := bson.M{
				"tags": bson.M{
					"$elemMatch": bson.M{
						"0": tagKey,
						"1": bson.M{"$in": tagValues},
					},
				},
			}
			tagQueries = append(tagQueries, qtf)
		}
		query = append(query, bson.E{Key: "$and", Value: tagQueries})
	}

	// Filter by Since (created_at >=)
	if fq.Since > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$gte": time.Unix(fq.Since, 0)}})
	}

	// Filter by Until (created_at <=)
	if fq.Until > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$lte": time.Unix(fq.Until, 0)}})
	}

	// Add Limit to options
	if fq.Limit > 0 {
		opts.SetLimit(int64(fq.Limit))
	}

	return query, opts, nil
}
