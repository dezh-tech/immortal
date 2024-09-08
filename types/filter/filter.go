package filter

import (
	"github.com/dezh-tech/immortal/types"
	"github.com/mailru/easyjson"
)

// Filter defined the filter structure based on NIP-01 and NIP-50.
type Filter struct {
	IDs     []string             `json:"ids"`
	Authors []string             `json:"authors"`
	Kinds   []types.Kind         `json:"kinds"`
	Tags    map[string]types.Tag `json:"tags"`
	Since   int64                `json:"since"`
	Until   int64                `json:"until"`
	Limit   uint16               `json:"limit"`

	// Sould we proxy Searchs to index server and elastic search?
	Search string `json:"search"` // Check NIP-50
}

// Decode decodes a byte array into event structure.
func Decode(b []byte) (*Filter, error) {
	e := new(Filter)

	if err := easyjson.Unmarshal(b, e); err != nil {
		return nil, err
	}

	return e, nil
}

// Encode encodes an event to a byte array.
func (e *Filter) Encode() ([]byte, error) {
	ee, err := easyjson.Marshal(e)
	if err != nil {
		return nil, err
	}

	return ee, nil
}
