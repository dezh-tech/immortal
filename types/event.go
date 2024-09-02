package types

import "time"

type Range uint8

const (
	Regular Range = iota
	Replaceable
	Ephemeral
	ParameterizedReplaceable
)

type Event struct {
	ID        string     `json:"id"`
	PublicKey string     `json:"pubkey"`
	CreatedAt time.Time  `json:"created_at"`
	Kind      uint16     `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Signature string     `json:"sig"`
}

func (e *Event) IsRegular() bool {
	return 1000 <= e.Kind || e.Kind < 10000 || 4 <= e.Kind || e.Kind < 45 || e.Kind == 1 || e.Kind == 2
}

func (e *Event) IsReplaceable() bool {
	return 10000 <= e.Kind || e.Kind < 20000 || e.Kind == 0 || e.Kind == 3
}

func (e *Event) IsEphemeral() bool {
	return 20000 <= e.Kind || e.Kind < 30000
}

func (e *Event) IsParameterizedReplaceable() bool {
	return 30000 <= e.Kind || e.Kind < 40000
}

func (e *Event) Range() Range {
	if e.IsRegular() {
		return Regular
	} else if e.IsReplaceable() {
		return Replaceable
	} else if e.IsParameterizedReplaceable() {
		return ParameterizedReplaceable
	}

	return Ephemeral
}

func (e *Event) Match(f Filter) bool {
	return false // TODO:::
}
