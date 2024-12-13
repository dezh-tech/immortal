package config

import (
	kraken "github.com/dezh-tech/immortal/client/gen/go"
	"github.com/dezh-tech/immortal/handler"
	"github.com/dezh-tech/immortal/server"
	"github.com/dezh-tech/immortal/utils"
)

type Parameters struct {
	Handler         *handler.Config
	WebsocketServer *server.Config
}

func (c *Config) LoadParameters(params *kraken.GetConfigResponse) error {
	url, err := utils.ParseURL(params.Url)
	if err != nil {
		return err
	}

	c.WebsocketServer.URL = url

	c.WebsocketServer.Limitation = &server.Limitation{
		MaxMessageLength:    params.Limitations.MaxMessageLength,
		MaxSubscriptions:    params.Limitations.MaxSubscriptions,
		MaxFilters:          params.Limitations.MaxFilters,
		MaxSubidLength:      params.Limitations.MaxSubidLength,
		MinPowDifficulty:    params.Limitations.MinPowDifficulty,
		AuthRequired:        params.Limitations.AuthRequired,
		PaymentRequired:     params.Limitations.PaymentRequired,
		RestrictedWrites:    params.Limitations.RestrictedWrites,
		MaxEventTags:        params.Limitations.MaxEventTags,
		MaxContentLength:    params.Limitations.MaxContentLength,
		CreatedAtLowerLimit: params.Limitations.CreatedAtLowerLimit,
		CreatedAtUpperLimit: params.Limitations.CreatedAtUpperLimit,
	}

	c.Handler = handler.Config{
		DefaultQueryLimit: params.Limitations.DefaultQueryLimit,
		MaxQueryLimit:     params.Limitations.MaxQueryLimit,
	}

	return nil
}
