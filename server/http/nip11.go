package http

import (
	"net/http"

	"github.com/dezh-tech/immortal/config"
)

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
	MaxMessageLength int  `json:"max_message_length,omitempty"`
	MaxSubscriptions int  `json:"max_subscriptions,omitempty"`
	MaxFilters       int  `json:"max_filters,omitempty"`
	MaxLimit         int  `json:"max_limit,omitempty"`
	MaxSubidLength   int  `json:"max_subid_length,omitempty"`
	MaxEventTags     int  `json:"max_event_tags,omitempty"`
	MaxContentLength int  `json:"max_content_length,omitempty"`
	MinPowDifficulty int  `json:"min_pow_difficulty,omitempty"`
	AuthRequired     bool `json:"auth_required"`
	PaymentRequired  bool `json:"payment_required"`
	RestrictedWrites bool `json:"restricted_writes"`
}

type RelayFeesDocument struct {
	Admission    []config.Admission
	Subscription []config.Subscription
	Publication  []config.Publication
}

// todo::: update with changes on db.
func fromConfig(cfg config.Config) *RelayInformationDocument {
	return &RelayInformationDocument{
		Name:          cfg.Parameters.Name,
		Description:   cfg.Parameters.Description,
		PubKey:        cfg.Parameters.Pubkey,
		SupportedNIPs: cfg.Parameters.SupportedNips,
		Software:      cfg.Parameters.Software,
		Version:       cfg.Parameters.Version,
		Contact:       cfg.Parameters.Contact,
		Limitation: &RelayLimitationDocument{
			MaxMessageLength: cfg.WebsocketServer.Limitation.MaxMessageLength,
			MaxSubscriptions: cfg.WebsocketServer.Limitation.MaxSubscriptions,
			MaxFilters:       cfg.WebsocketServer.Limitation.MaxFilters,
			MaxLimit:         cfg.Parameters.Handler.Limitation.MaxLimit,
			MaxSubidLength:   cfg.WebsocketServer.Limitation.MaxSubidLength,
			MaxEventTags:     cfg.Parameters.Handler.Limitation.MaxEventTags,
			MaxContentLength: cfg.Parameters.Handler.Limitation.MaxContentLength,
			MinPowDifficulty: cfg.WebsocketServer.Limitation.MinPowDifficulty,
			AuthRequired:     cfg.WebsocketServer.Limitation.AuthRequired,
			PaymentRequired:  cfg.WebsocketServer.Limitation.PaymentRequired,
			RestrictedWrites: cfg.WebsocketServer.Limitation.RestrictedWrites,
		},
		RelayCountries: cfg.Parameters.RelayCountries,
		LanguageTags:   cfg.Parameters.LanguageTags,
		Tags:           cfg.Parameters.Tags,
		PostingPolicy:  cfg.Parameters.PostingPolicy,
		PaymentsURL:    cfg.Parameters.PaymentsURL,
		Icon:           cfg.Parameters.Icon,
		Fees: &RelayFeesDocument{
			Admission:    cfg.Parameters.Fees.Admission,
			Subscription: cfg.Parameters.Fees.Subscription,
			Publication:  cfg.Parameters.Fees.Publication,
		},
	}
}

func (s *Server) nip11Handler(w http.ResponseWriter, r *http.Request) {
	nip11Doc := fromConfig(*s.config)
	s.respondWithJSON(w, 200, nip11Doc)
}
