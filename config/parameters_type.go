package config

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
