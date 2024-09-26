package websocket

type Limitation struct {
	MaxMessageLength int  `bson:"max_message_length" json:"max_message_length"`
	MaxSubscriptions int  `bson:"max_subscriptions"  json:"max_subscriptions"`
	MaxFilters       int  `bson:"max_filters"        json:"max_filters"`
	MaxSubidLength   int  `bson:"max_subid_length"   json:"max_subid_length"`
	MinPowDifficulty int  `bson:"min_pow_difficulty" json:"min_pow_difficulty"`
	AuthRequired     bool `bson:"auth_required"      json:"auth_required"`
	PaymentRequired  bool `bson:"payment_required"   json:"payment_required"`
	RestrictedWrites bool `bson:"restricted_writes"  json:"restricted_writes"`
}
type Config struct {
	Bind            string `yaml:"bind"`
	BloomBackupPath string `yaml:"bloom_backup_path"`
	Port            uint16 `yaml:"port"`
	KnownBloomSize  uint   `yaml:"known_bloom_size"`
	Limitation      *Limitation
}
