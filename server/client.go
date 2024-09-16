package server

import (
	"sync"

	"github.com/dezh-tech/immortal/types/filter"
)

type client struct {
	subs map[string]filter.Filters
	*sync.RWMutex
}
