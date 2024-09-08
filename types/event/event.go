package event

import (
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/filter"
	"github.com/mailru/easyjson"
)

// Event represents an event structure defined on NIP-01.
type Event struct {
	ID        string      `json:"id"`
	PublicKey string      `json:"pubkey"`
	CreatedAt int64       `json:"created_at"`
	Kind      types.Kind  `json:"kind"`
	Tags      []types.Tag `json:"tags"`
	Content   string      `json:"content"`
	Signature string      `json:"sig"`
}

// Match checks if the event is match with given filter.
// Note: this method intended to be used for already open subscriptions and recently received events.
// For new subscriptions and queries for stored data use the database query and don't use this to verify the result.
func (e *Event) Match(f filter.Filter) bool {
	if e == nil {
		return false
	}

	if f.IDs != nil && !types.ContainsString(e.ID, f.IDs) {
		return false
	}

	if f.Authors != nil && !types.ContainsString(e.PublicKey, f.Authors) {
		return false
	}

	if f.Kinds != nil && !types.ContainsKind(e.Kind, f.Kinds) {
		return false
	}

	if e.CreatedAt >= f.Since || e.CreatedAt <= f.Until {
		return false
	}

	for f, vals := range f.Tags {
		for _, t := range e.Tags {
			if len(t) < 2 {
				continue
			}

			if f != "#"+t[0] { // TODO:: should we replace + with strings.Builder?
				return false
			}

			var containsValue bool
			for _, v := range vals {
				if v == t[1] {
					containsValue = true
					break
				}
			}

			if !containsValue {
				return false
			}
		}
	}

	return true
}

// Decode decodes a byte array into Event structure.
func Decode(b []byte) (*Event, error) {
	e := new(Event)

	if err := easyjson.Unmarshal(b, e); err != nil {
		return nil, err
	}

	return e, nil
}

// Encode encodes an Event to a byte array.
func (e *Event) Encode() ([]byte, error) {
	ee, err := easyjson.Marshal(e)
	if err != nil {
		return nil, err
	}

	return ee, nil
}

// IsValid function validats an event Signature and ID.
func (e *Event) IsValid() bool {
	return false // TODO:::
}
