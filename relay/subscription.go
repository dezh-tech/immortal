package relay

import "github.com/dezh-tech/immortal/types/filter"

type Subscription struct {
	ID      string
	Filters filter.Filters
}
