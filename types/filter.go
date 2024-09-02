package types

import "time"

type Filter struct {
	IDs     []string  `json:"ids"`
	Authors []string  `json:"authors"`
	Kinds   []uint16  `json:"kinds"`
	Tags    []string  `json:"tags"` // Is that correct? // TODO:::
	Since   time.Time `json:"since"`
	Until   time.Time `json:"until"`
	Limit   int16     `json:"limit"`
}
