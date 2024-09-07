package types

// Filter defined the filter structure based on NIP-01 and NIP-50.
type Filter struct {
	IDs     []string            `json:"ids"`
	Authors []string            `json:"authors"`
	Kinds   []uint16            `json:"kinds"`
	Tags    map[string][]string `json:"tags"`
	Since   int64               `json:"since"`
	Until   int64               `json:"until"`
	Limit   int16               `json:"limit"`

	// Sould we proxy Searchs to index server and elastic search?
	Search string `json:"search"` // Check NIP-50
}
