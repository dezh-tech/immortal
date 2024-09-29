package handler

type Limitation struct {
	MaxLimit uint16 `bson:"max_limit" json:"max_limit"`
}

type Config struct {
	InitialQueryDefaultLimit int64      `bson:"initial_query_default_limit"`
	Limitation               Limitation `bson:"limitation"`
}
