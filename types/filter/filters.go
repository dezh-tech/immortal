package filter

import (
	"encoding/json"

	"github.com/dezh-tech/immortal/types/event"
)

type Filters []Filter

// String returns and string representation of encoded filters.
func (f Filters) String() string {
	j, err := json.Marshal(f)
	if err != nil {
		return ""
	}

	return string(j)
}

// Match checks id the given event e is match with any of filters f.
func (f Filters) Match(e *event.Event) bool {
	for _, filter := range f {
		if filter.Match(e) {
			return true
		}
	}

	return false
}
