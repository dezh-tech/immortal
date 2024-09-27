package handler

type Limitation struct {
	MaxLimit uint16 `bson:"max_limit"              json:"max_limit"`
	// todo(@zig)::: move to server.
	MaxEventTags uint `bson:"max_event_tags"         json:"max_event_tags"`
	// todo(@zig)::: move to server.
	MaxContentLength uint `bson:"max_content_length"     json:"max_content_length"`
	// todo(@zig)::: move to server.
	CreatedAtLowerLimit uint `bson:"created_at_lower_limit" json:"created_at_lower_limit"`
	// todo(@zig)::: move to server.
	CreatedAtUpperLimit uint `bson:"created_at_upper_limit" json:"created_at_upper_limit"`
}

type Config struct {
	InitialQueryDefaultLimit int64      `bson:"initial_query_default_limit"`
	Limitation               Limitation `bson:"limitation"`
}
