package filter

import (
	"github.com/dezh-tech/immortal/types"
	"github.com/dezh-tech/immortal/types/event"
	"github.com/mailru/easyjson"
)

// Filter defined the filter structure based on NIP-01 and NIP-50.
type Filter struct {
	IDs     []string     `json:"ids"`
	Authors []string     `json:"authors"`
	Kinds   []types.Kind `json:"kinds"`
	Tags    map[string][]string
	Since   int64  `json:"since"`
	Until   int64  `json:"until"`
	Limit   uint32 `json:"limit"`
	Search  string `json:"search"` // see NIP-50.
}

// todo::: can/should we support nip-50 in match function?
// Match checks if the event is match with given filter.
// Note: this method intended to be used for already open subscriptions and recently received events.
// For new subscriptions and queries for stored data use the database query and don't use this to verify the result.
func (f *Filter) Match(e *event.Event, pubkey string) bool {
	if e == nil {
		return false
	}

	if len(f.IDs) != 0 && !types.ContainsString(e.ID, f.IDs) {
		return false
	}

	if len(f.Authors) != 0 && !types.ContainsString(e.PublicKey, f.Authors) {
		return false
	}

	if len(f.Kinds) != 0 && !types.ContainsKind(e.Kind, f.Kinds) {
		return false
	}

	if f.Since != 0 && e.CreatedAt < f.Since {
		return false
	}

	if f.Until != 0 && e.CreatedAt > f.Until {
		return false
	}

	for tagName, vals := range f.Tags {
		if vals != nil && !e.Tags.ContainsAny(tagName, vals) {
			return false
		}
	}

	if e.Kind == types.KindGiftWrap && !e.Tags.ContainsTag("p", pubkey) {
		return false
	}

	return true
}

// Decode decodes a byte array into event structure.
func Decode(b []byte) (*Filter, error) {
	e := new(Filter)

	if err := easyjson.Unmarshal(b, e); err != nil {
		return nil, types.DecodeError{
			Reason: err.Error(),
		}
	}

	return e, nil
}

// Encode encodes an event to a byte array.
func (f *Filter) Encode() ([]byte, error) {
	ee, err := easyjson.Marshal(f)
	if err != nil {
		return nil, types.EncodeError{
			Reason: err.Error(),
		}
	}

	return ee, nil
}

// String returns and string representation of encoded filter.
func (f *Filter) String() string {
	ef, err := f.Encode()
	if err != nil {
		return ""
	}

	return string(ef)
}
