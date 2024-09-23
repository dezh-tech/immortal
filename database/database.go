package database

import (
	"context"
	"time"

	"github.com/dezh-tech/immortal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var KindToCollectionName = map[types.Kind]string{
	types.KindTextNote:        "text_notes",
	types.KindReaction:        "reactions",
	types.KindProfileMetadata: "profile_metadatas",
}

type Database struct {
	DBName       string
	QueryTimeout time.Duration
	Client       *mongo.Client
}

func New(cfg Config) (*Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.URI).
		SetServerAPIOptions(serverAPI).
		SetConnectTimeout(time.Duration(cfg.ConnectionTimeout) * time.Millisecond). // Convert to time.Duration
		SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilSliceAsEmpty:   true,
		})

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).
		Decode(&result); err != nil {
		return nil, err
	}

	return &Database{
		Client:       client,
		DBName:       cfg.DBName,
		QueryTimeout: time.Duration(cfg.QueryTimeout) * time.Millisecond,
	}, nil
}

func (db *Database) Stop() error {
	return db.Client.Disconnect(context.TODO())
}
