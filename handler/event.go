package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dezh-tech/immortal/types/event"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) HandleEvent(e *event.Event) error {
	collName, ok := KindToCollectionName[e.Kind]
	if !ok {
		return fmt.Errorf("kind %d is not supported yet", e.Kind)
	}

	coll := h.DB.Client.Database(h.DB.DBName).Collection(collName)
	ctx, cancel := context.WithTimeout(context.Background(), h.DB.QueryTimeout)
	defer cancel()

	if e.Kind.IsRegular() {
		_, err := coll.InsertOne(ctx, e)
		if err != nil {
			return err
		}

		return nil
	}

	var query bson.D
	var filter bson.D

	if e.Kind.IsReplaceable() {
		query = bson.D{
			bson.E{
				Key:   "pubkey",
				Value: e.PublicKey,
			},
			bson.E{
				Key:   "kind",
				Value: e.Kind,
			},
		}

		filter = bson.D{
			bson.E{
				Key:   "pubkey",
				Value: e.PublicKey,
			},
			{
				Key:   "kind",
				Value: e.Kind,
			},
			{
				Key: "created_at",
				Value: bson.M{
					"$lte": time.Unix(e.CreatedAt, 0),
				},
			},
		}
	}

	if e.Kind.IsAddressable() {
		var dTag string
		for _, t := range e.Tags {
			if len(t) < 2 {
				continue
			}

			if t[0] == "d" {
				dTag = t[1]

				break
			}
		}

		if dTag == "" {
			return errors.New("no d tag found")
		}

		query = bson.D{
			bson.E{
				Key:   "pubkey",
				Value: e.PublicKey,
			},
			bson.E{
				Key:   "kind",
				Value: e.Kind,
			},
			bson.E{
				Key: "$and",
				Value: bson.M{
					"tags": bson.M{
						"$elemMatch": bson.M{
							"0": "d",
							"1": dTag,
						},
					},
				},
			},
		}

		filter = bson.D{
			bson.E{
				Key:   "pubkey",
				Value: e.PublicKey,
			},
			{
				Key:   "kind",
				Value: e.Kind,
			},
			{
				Key: "created_at",
				Value: bson.M{
					"$lte": time.Unix(e.CreatedAt, 0),
				},
			},
			bson.E{
				Key: "$and",
				Value: bson.M{
					"tags": bson.M{
						"$elemMatch": bson.M{
							"0": "d",
							"1": dTag,
						},
					},
				},
			},
		}
	}

	cursor, err := coll.Find(ctx, query)
	if err != nil {
		return err
	}

	var result []event.Event
	if err := cursor.All(ctx, &result); err != nil {
		return err
	}

	if len(result) == 0 {
		_, err := coll.InsertOne(ctx, e)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = coll.ReplaceOne(ctx, filter, e, &options.ReplaceOptions{})
	if err != nil {
		return err
	}

	return nil
}
