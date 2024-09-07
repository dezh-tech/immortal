package types

// Range is the events kind ranges.
type Range uint8

type Tag [2]string

const (
	Regular Range = iota
	Replaceable
	Ephemeral
	ParameterizedReplaceable
)

// Event reperesents an event structure defined on NIP-01.
type Event struct {
	ID        string `json:"id"`
	PublicKey string `json:"pubkey"`
	CreatedAt int64  `json:"created_at"`
	Kind      uint16 `json:"kind"`
	Tags      []Tag  `json:"tags"`
	Content   string `json:"content"`
	Signature string `json:"sig"`
}

// IsRegular checks if the gived event kind is in Regular range.
func (e *Event) IsRegular() bool {
	return 1000 <= e.Kind || e.Kind < 10000 || 4 <= e.Kind || e.Kind < 45 || e.Kind == 1 || e.Kind == 2
}

// IsReplaceable checks if the gived event kind is in Replaceable range.
func (e *Event) IsReplaceable() bool {
	return 10000 <= e.Kind || e.Kind < 20000 || e.Kind == 0 || e.Kind == 3
}

// IsEphemeral checks if the gived event kind is in Ephemeral range.
func (e *Event) IsEphemeral() bool {
	return 20000 <= e.Kind || e.Kind < 30000
}

// IsParameterizedReplaceable checks if the gived event kind is in ParameterizedReplaceable range.
func (e *Event) IsParameterizedReplaceable() bool {
	return 30000 <= e.Kind || e.Kind < 40000
}

// Range returns the events kind range based on NIP-01.
func (e *Event) Range() Range {
	if e.IsRegular() {
		return Regular
	} else if e.IsReplaceable() {
		return Replaceable
	} else if e.IsParameterizedReplaceable() {
		return ParameterizedReplaceable
	}

	return Ephemeral
}

// Match checks if the event is match with given filter.
// Note: this method intended to be used for already open subscriptions and recently received events.
// For new subscriptions and queries for stored data use the database query and don't use this to verify the result.
func (e *Event) Match(f Filter) bool {
	if e == nil {
		return false
	}

	if f.IDs != nil && !ContainsString(e.ID, f.IDs) {
		return false
	}

	if f.Authors != nil && !ContainsString(e.PublicKey, f.Authors) {
		return false
	}

	if f.Kinds != nil && !ContainsUint16(e.Kind, f.Kinds) {
		return false
	}

	for f, vals := range f.Tags {
		for _, t := range e.Tags {
			if f != "#"+t[0] { // should we change it(+)?
				return false
			}

			var containsValue bool
			for _, v := range vals {
				if v == t[1] {
					containsValue = true
					break
				}
			}

			if !containsValue {
				return false
			}
		}
	}

	return true
}

// Decode decodes a byte array into event structure.
func Decode(b []byte) (*Event, error) {
	return nil, nil // TODO:::
}

// Encode encodes an event to a byte array.
func (e *Event) Encode() ([]byte, error) {
	return nil, nil // TODO:::
}

// IsValid function validats an event Signature and ID.
func (e *Event) IsValid() bool {
	return false // TODO:::
}
