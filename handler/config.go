package handler

type Limitation struct {
	MaxLimit            int `bson:"max_limit"              json:"max_limit"`
	MaxEventTags        int `bson:"max_event_tags"         json:"max_event_tags"`
	MaxContentLength    int `bson:"max_content_length"     json:"max_content_length"`
	CreatedAtLowerLimit int `bson:"created_at_lower_limit" json:"created_at_lower_limit"`
	CreatedAtUpperLimit int `bson:"created_at_upper_limit" json:"created_at_upper_limit"`
}

type Config struct {
	InitialQueryDefaultLimit int64      `bson:"initial_query_default_limit"`
	Limitation               Limitation `bson:"limitation"`
}
