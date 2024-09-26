package websocket

import (
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
)

type clientState struct {
	subs map[string]filter.Filters
	*sync.RWMutex
}
