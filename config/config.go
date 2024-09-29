package config

import (
	"os"

	"github.com/dezh-tech/immortal/database"
	"github.com/dezh-tech/immortal/server/http"
	"github.com/dezh-tech/immortal/server/websocket"
	"github.com/dezh-tech/immortal/types/nip11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config reprsents the configs used by relay and other concepts on system.
type Config struct {
	Environment     string           `yaml:"environment"`
	WebsocketServer websocket.Config `yaml:"ws_server"`
	HTTPServer      http.Config      `yaml:"http_server"`
	Database        database.Config  `yaml:"database"`
	Parameters      *Parameters
}

// Load loads config from file and env.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}
	defer file.Close()

	config := &Config{}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(config); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	if config.Environment != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, Error{
				reason: err.Error(),
			}
		}
	}

	config.Database.URI = os.Getenv("IMMO_MONGO_URI")

	if err = config.basicCheck(); err != nil {
		return nil, Error{
			reason: err.Error(),
		}
	}

	return config, nil
}

func (c *Config) GetNIP11Documents() *nip11.RelayInformationDocument {
	n11d := &nip11.RelayInformationDocument{
		Name:          c.Parameters.Name,
		Description:   c.Parameters.Description,
		PubKey:        c.Parameters.Pubkey,
		SupportedNIPs: c.Parameters.SupportedNips,
		Software:      c.Parameters.Software,
		Version:       c.Parameters.Version,
		Contact:       c.Parameters.Contact,
		Limitation: &nip11.RelayLimitationDocument{
			MaxMessageLength: c.WebsocketServer.Limitation.MaxMessageLength,
			MaxSubscriptions: c.WebsocketServer.Limitation.MaxSubscriptions,
			MaxFilters:       c.WebsocketServer.Limitation.MaxFilters,
			MaxLimit:         c.Parameters.Handler.Limitation.MaxLimit,
			MaxSubidLength:   c.WebsocketServer.Limitation.MaxSubidLength,
			MaxEventTags:     c.Parameters.Handler.Limitation.MaxEventTags,
			MaxContentLength: c.Parameters.Handler.Limitation.MaxContentLength,
			MinPowDifficulty: c.WebsocketServer.Limitation.MinPowDifficulty,
			AuthRequired:     c.WebsocketServer.Limitation.AuthRequired,
			PaymentRequired:  c.WebsocketServer.Limitation.PaymentRequired,
			RestrictedWrites: c.WebsocketServer.Limitation.RestrictedWrites,
		},
		RelayCountries: c.Parameters.RelayCountries,
		LanguageTags:   c.Parameters.LanguageTags,
		Tags:           c.Parameters.Tags,
		PostingPolicy:  c.Parameters.PostingPolicy,
		PaymentsURL:    c.Parameters.PaymentsURL,
		Icon:           c.Parameters.Icon,
		Fees:           new(nip11.RelayFeesDocument),
	}

	addmissions := make([]nip11.Admission, 0)
	for _, a := range c.Parameters.Fees.Admission {
		addmissions = append(addmissions, nip11.Admission{
			Amount: a.Amount,
			Unit:   a.Unit,
		})
	}

	subscription := make([]nip11.Subscription, 0)
	for _, s := range c.Parameters.Fees.Subscription {
		subscription = append(subscription, nip11.Subscription{
			Amount: s.Amount,
			Unit:   s.Unit,
			Period: s.Period,
		})
	}

	publication := make([]nip11.Publication, 0)
	for _, p := range c.Parameters.Fees.Publication {
		publication = append(publication, nip11.Publication{
			Kinds:  p.Kinds,
			Amount: p.Amount,
			Unit:   p.Unit,
		})
	}

	n11d.Fees.Admission = addmissions
	n11d.Fees.Subscription = subscription
	n11d.Fees.Publication = publication

	return n11d
}

// basicCheck validates the basic stuff in config.
func (c *Config) basicCheck() error {
	return nil
}
