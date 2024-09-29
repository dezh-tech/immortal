package nip11

type RelayInformationDocument struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PubKey        string `json:"pubkey"`
	Contact       string `json:"contact"`
	SupportedNIPs []int  `json:"supported_nips"`
	Software      string `json:"software"`
	Version       string `json:"version"`

	Limitation     *RelayLimitationDocument `json:"limitation,omitempty"`
	RelayCountries []string                 `json:"relay_countries,omitempty"`
	LanguageTags   []string                 `json:"language_tags,omitempty"`
	Tags           []string                 `json:"tags,omitempty"`
	PostingPolicy  string                   `json:"posting_policy,omitempty"`
	PaymentsURL    string                   `json:"payments_url,omitempty"`
	Fees           *RelayFeesDocument       `json:"fees,omitempty"`
	Icon           string                   `json:"icon"`
}

type RelayLimitationDocument struct {
	MaxMessageLength int    `json:"max_message_length,omitempty"`
	MaxSubscriptions int    `json:"max_subscriptions,omitempty"`
	MaxFilters       int    `json:"max_filters,omitempty"`
	MaxLimit         uint16 `json:"max_limit,omitempty"`
	MaxSubidLength   int    `json:"max_subid_length,omitempty"`
	MaxEventTags     int    `json:"max_event_tags,omitempty"`
	MaxContentLength int    `json:"max_content_length,omitempty"`
	MinPowDifficulty int    `json:"min_pow_difficulty,omitempty"`
	AuthRequired     bool   `json:"auth_required"`
	PaymentRequired  bool   `json:"payment_required"`
	RestrictedWrites bool   `json:"restricted_writes"`
}

type RelayFeesDocument struct {
	Admission    []Admission
	Subscription []Subscription
	Publication  []Publication
}

type Subscription struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
	Period int    `json:"period"`
}

type Admission struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

type Publication struct {
	Kinds  []int  `json:"kinds"`
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

type Fees struct {
	Subscription []Subscription `json:"subscription,omitempty"`
	Publication  []Publication  `json:"publication,omitempty"`
	Admission    []Admission    `json:"admission,omitempty"`
}
