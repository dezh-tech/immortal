package config

import (
	"context"
	"errors"

	"github.com/dezh-tech/immortal"
	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Retention struct {
	Time  int         `bson:"time,omitempty"  json:"time,omitempty"`
	Count int         `bson:"count,omitempty" json:"count,omitempty"`
	Kinds interface{} `bson:"kinds,omitempty" json:"kinds,omitempty"`
}

type Fees struct {
	Amount int    `bson:"amount" json:"amount"`
	Unit   string `bson:"unit"   json:"unit"`
	Period int    `bson:"period" json:"period"`
}

type Subscription struct {
	Subscription []Fees `bson:"subscription" json:"subscription"`
}
type Parameters struct {
	Handler        *handler.Config `bson:"handler"         json:"handler"`
	Retention      *Retention      `bson:"retention"       json:"retention"`
	Fees           *Subscription   `bson:"fees"            json:"fees"`
	Name           string          `bson:"name"            json:"name"`
	Description    string          `bson:"description"     json:"description"`
	Pubkey         string          `bson:"pubkey"          json:"pubkey"`
	Software       string          `bson:"software"        json:"software"`
	SupportedNips  []int           `bson:"supported_nips"  json:"supported_nips"`
	Version        string          `bson:"version"         json:"version"`
	RelayCountries []string        `bson:"relay_countries" json:"relay_countries"`
	LanguageTags   []string        `bson:"language_tags"   json:"language_tags"`
	Tags           []string        `bson:"tags"            json:"tags"`
	PostingPolicy  string          `bson:"posting_policy"  json:"posting_policy"`
	PaymentsURL    string          `bson:"payments_url"    json:"payments_url"`
	Icon           string          `bson:"icon"            json:"icon"`
}

func (c *Config) LoadParameters(db *database.Database) error {
	coll := db.Client.Database(db.DBName).Collection("config_parameters")

	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	filter := bson.M{}

	var result *Parameters
	err := coll.FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// insert default parameters
		newDocument := &Parameters{
			Name:           "Immortal",                                                         // relay name
			Description:    "A Nostr relay for scale",                                          // description
			Pubkey:         "0000000000000000000000000000000000000000000000000000000000000000", // pubkey
			Software:       "https://github.com/dezh-tech/immortal",                            // software repository URL
			SupportedNips:  []int{1, 11},                                                       // Supported NIPs (protocols)
			Version:        immortal.StringVersion(),                                           // Version of the relay software
			RelayCountries: []string{"US"},                                                     // country support
			LanguageTags:   []string{"en"},                                                     // language tags
			Tags:           []string{},                                                         // tags
			PostingPolicy:  "",                                                                 // posting policy URL
			PaymentsURL:    "",                                                                 // payments URL
			Icon:           "",                                                                 // icon URL
			Handler: &handler.Config{
				InitialQueryDefaultLimit: 100,
				Limitation: handler.Limitation{
					MaxMessageLength:    8192,  // Maximum length of a single message (in bytes or characters)
					MaxSubscriptions:    20,    // Maximum number of concurrent subscriptions a client can create
					MaxFilters:          20,    // Maximum number of filters a client can apply in a subscription
					MaxLimit:            1000,  // Maximum number of events returned in a query
					MaxSubidLength:      256,   // Maximum length of a subscription identifier
					MaxEventTags:        200,   // Maximum number of tags allowed in a single event
					MaxContentLength:    4096,  // Maximum content length of an event (in bytes)
					MinPowDifficulty:    0,     // Minimum proof-of-work difficulty for publishing events
					AuthRequired:        false, // Whether authentication is required for writes
					PaymentRequired:     false, // Whether payment is required to interact with the relay
					RestrictedWrites:    false, // Whether writes are restricted to authenticated or paying users
					CreatedAtLowerLimit: 0,     // Earliest timestamp allowed for event creation
					CreatedAtUpperLimit: 0,     // Latest timestamp allowed for event creation (0 for no limit)
				},
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

func (c *Config) SetParameters(db *database.Database, params *Parameters) error {
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
