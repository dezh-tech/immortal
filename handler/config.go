package handler

type Limitation struct {
	MaxMessageLength    int  `bson:"max_message_length"     json:"max_message_length"`
	MaxSubscriptions    int  `bson:"max_subscriptions"      json:"max_subscriptions"`
	MaxFilters          int  `bson:"max_filters"            json:"max_filters"`
	MaxLimit            int  `bson:"max_limit"              json:"max_limit"`
	MaxSubidLength      int  `bson:"max_subid_length"       json:"max_subid_length"`
	MaxEventTags        int  `bson:"max_event_tags"         json:"max_event_tags"`
	MaxContentLength    int  `bson:"max_content_length"     json:"max_content_length"`
	MinPowDifficulty    int  `bson:"min_pow_difficulty"     json:"min_pow_difficulty"`
	AuthRequired        bool `bson:"auth_required"          json:"auth_required"`
	PaymentRequired     bool `bson:"payment_required"       json:"payment_required"`
	RestrictedWrites    bool `bson:"restricted_writes"      json:"restricted_writes"`
	CreatedAtLowerLimit int  `bson:"created_at_lower_limit" json:"created_at_lower_limit"`
	CreatedAtUpperLimit int  `bson:"created_at_upper_limit" json:"created_at_upper_limit"`
}

type Config struct {
	InitialQueryDefaultLimit int64      `bson:"initial_query_default_limit"`
	Limitation               Limitation `bson:"limitation"`
}
