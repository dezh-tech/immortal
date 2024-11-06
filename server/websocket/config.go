package websocket

type Limitation struct {
	MaxMessageLength    int   `bson:"max_message_length"     json:"max_message_length"` // todo?.
	MaxSubscriptions    int   `bson:"max_subscriptions"      json:"max_subscriptions"`
	MaxFilters          int   `bson:"max_filters"            json:"max_filters"`
	MaxSubidLength      int   `bson:"max_subid_length"       json:"max_subid_length"`
	MinPowDifficulty    int   `bson:"min_pow_difficulty"     json:"min_pow_difficulty"`
	AuthRequired        bool  `bson:"auth_required"          json:"auth_required"`
	PaymentRequired     bool  `bson:"payment_required"       json:"payment_required"`  // todo.
	RestrictedWrites    bool  `bson:"restricted_writes"      json:"restricted_writes"` // todo.
	MaxEventTags        int   `bson:"max_event_tags"         json:"max_event_tags"`
	MaxContentLength    int   `bson:"max_content_length"     json:"max_content_length"`
	CreatedAtLowerLimit int64 `bson:"created_at_lower_limit" json:"created_at_lower_limit"`
	CreatedAtUpperLimit int64 `bson:"created_at_upper_limit" json:"created_at_upper_limit"`
}

type Config struct {
	Bind       string `yaml:"bind"`
	Port       uint16 `yaml:"port"`
	Limitation *Limitation
}
