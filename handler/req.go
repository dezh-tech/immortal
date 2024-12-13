package handler

import (
	"context"

	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/dezh-tech/immortal/types/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var possibleKinds = []types.Kind{
	types.KindUserMetadata,
	types.KindShortTextNote,
	types.KindZap,
	types.KindRelayListMetadata,
}

type filterQuery struct {
	Tags map[string][]string

	Authors []string
	IDs     []string

	Since int64
	Until int64
	Limit uint32
}

func (h *Handler) HandleReq(fs filter.Filters) ([]event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.db.QueryTimeout)
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

		if len(f.Kinds) != 0 {
			uniqueKinds := removeDuplicateKind(f.Kinds)
			for _, k := range uniqueKinds {
				queryKinds[k] = append(queryKinds[k], qf)
			}
		} else {
			//! we query most requested kinds if there is no kind provided.
			// FIX: any better way?
			for _, k := range possibleKinds {
				queryKinds[k] = append(queryKinds[k], qf)
			}
		}
	}

	var finalResult []event.Event

	for kind, filters := range queryKinds {
		collection := h.db.Client.Database(h.db.DBName).Collection(getCollectionName(kind))
		for _, f := range filters {
			query, opts, err := h.FilterToQuery(&f)
			if err != nil {
				continue
			}

			cursor, err := collection.Find(ctx, query, opts)
			if err != nil {
				return nil, err
			}

			var result []event.Event
			if err := cursor.All(ctx, &result); err != nil {
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

func (h *Handler) FilterToQuery(fq *filterQuery) (bson.D, *options.FindOptions, error) {
	query := make(bson.D, 0)
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
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$gte": fq.Since}})
	}

	// Filter by Until (created_at <=)
	if fq.Until > 0 {
		query = append(query, bson.E{Key: "created_at", Value: bson.M{"$lte": fq.Since}})
	}

	// Add Limit to options
	if fq.Limit > 0 && fq.Limit < h.config.MaxQueryLimit {
		opts.SetLimit(int64(fq.Limit))
	} else {
		opts.SetLimit(int64(h.config.DefaultQueryLimit))
	}

	opts.SetSort(bson.D{
		{Key: "created_at", Value: -1},
		{Key: "id", Value: 1},
	})

	return query, opts, nil
}
