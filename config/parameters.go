package config

import (
	"github.com/dezh-tech/immortal/delivery/websocket"
	mpb "github.com/dezh-tech/immortal/infrastructure/grpc_client/gen"
	"github.com/dezh-tech/immortal/repository"
	"github.com/dezh-tech/immortal/utils"
)

type Parameters struct {
	Handler         *repository.Config
	WebsocketServer *websocket.Config
}

func (c *Config) LoadParameters(params *mpb.GetParametersResponse) error {
	url, err := utils.ParseURL(params.Url)
	if err != nil {
		return err
	}

	c.WebsocketServer.URL = url

	c.WebsocketServer.Limitation = &websocket.Limitation{
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

	c.Handler = repository.Config{
		DefaultQueryLimit: params.Limitations.DefaultQueryLimit,
		MaxQueryLimit:     params.Limitations.MaxQueryLimit,
	}

	return nil
}
