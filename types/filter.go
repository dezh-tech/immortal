package types

type Filter struct {
	IDs     []string            `json:"ids"`
	Authors []string            `json:"authors"`
	Kinds   []uint16            `json:"kinds"`
	Tags    map[string][]string `json:"tags"` // Is that correct? // TODO:::
	Since   int64               `json:"since"`
	Until   int64               `json:"until"`
	Limit   int16               `json:"limit"`
}

func (f *Filter) IsValid() bool {
	return false // TODO::
}
