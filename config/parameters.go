package config

import (
	"context"
	"errors"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Parameters struct {
	Handler *handler.Config `bson:"handler" `
}

func (c *Config) LoadParameters(db *database.Database) error {
	coll := db.Client.Database(db.DBName).Collection("config_parameters")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{}

	var result Parameters
	err := coll.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// insert default parameters
		newDocument := Parameters{
			Handler: &handler.Config{
				InitialQueryDefaultLimit: 100,
			},
		}

		insertErr := c.SetParameters(db, newDocument)
		if insertErr != nil {
			return insertErr
		}

		return nil
	} else if err != nil {
		return err
	}

	c.Parameters = result

	return nil
}

func (c *Config) SetParameters(db *database.Database, params Parameters) error {
	coll := db.Client.Database(db.DBName).Collection("config_parameters")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()
	_, insertErr := coll.InsertOne(ctx, params)
	if insertErr != nil {
		return insertErr
	}
	c.Parameters = params

	return nil
}
