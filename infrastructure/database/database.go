package database

import (
	"context"
	"log"
	"time"

	"github.com/dezh-tech/immortal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DBName       string
	QueryTimeout time.Duration
	Client       *mongo.Client
}

func Connect(cfg Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnectionTimeout)*time.Millisecond)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.URI).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(time.Duration(cfg.ConnectionTimeout) * time.Millisecond).
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilSliceAsEmpty:   true,
		})

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.QueryTimeout)*time.Millisecond)
	defer cancel()

	if err := client.Ping(qCtx, nil); err != nil {
		return nil, err
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("uq_id"),
	}

	for _, collName := range types.KindToName {
		qCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.QueryTimeout)*time.Millisecond)
		_, err := client.Database(cfg.DBName).Collection(collName).Indexes().CreateOne(qCtx, indexModel)
		if err != nil {
			cancel()

			continue
		}
		cancel()
	}

	return &Database{
		Client:       client,
		DBName:       cfg.DBName,
		QueryTimeout: time.Duration(cfg.QueryTimeout) * time.Millisecond,
	}, nil
}

func (db *Database) Stop() error {
	log.Println("closing database connection...")

	return db.Client.Disconnect(context.Background())
}
