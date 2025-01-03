package websocket

import (
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
)

type clientState struct {
	challenge string
	pubkey    *string
	isKnown   *bool
	subs      map[string]filter.Filters
	*sync.RWMutex
}
