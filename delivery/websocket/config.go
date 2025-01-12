package websocket

import "net/url"

type Limitation struct {
	MaxMessageLength    int32 // todo?.
	MaxSubscriptions    int32
	MaxSubidLength      int32
	MinPowDifficulty    int32
	AuthRequired        bool
	PaymentRequired     bool
	RestrictedWrites    bool
	MaxEventTags        int32
	MaxContentLength    int32
	CreatedAtLowerLimit int64
	CreatedAtUpperLimit int64
}

type Config struct {
	Bind       string `yaml:"bind"`
	Port       uint16 `yaml:"port"`
	URL        *url.URL
	Limitation *Limitation
}
