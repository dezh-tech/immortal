package config

import (
	"github.com/dezh-tech/immortal/delivery/websocket"
	mpb "github.com/dezh-tech/immortal/infrastructure/grpc_client/gen"
	"github.com/dezh-tech/immortal/pkg/utils"
)

func (c *Config) LoadParameters(params *mpb.GetParametersResponse) error {
	url, err := utils.ParseURL(params.Url)
	if err != nil {
		return err
	}

	c.WebsocketServer.SetURL(url)

	c.WebsocketServer.SetLimitation(
		&websocket.Limitation{
			MaxMessageLength:    params.Limitations.MaxMessageLength,
			MaxSubscriptions:    params.Limitations.MaxSubscriptions,
			MaxSubidLength:      params.Limitations.MaxSubidLength,
			MinPowDifficulty:    params.Limitations.MinPowDifficulty,
			AuthRequired:        params.Limitations.AuthRequired,
			PaymentRequired:     params.Limitations.PaymentRequired,
			RestrictedWrites:    params.Limitations.RestrictedWrites,
			MaxEventTags:        params.Limitations.MaxEventTags,
			MaxContentLength:    params.Limitations.MaxContentLength,
			CreatedAtLowerLimit: params.Limitations.CreatedAtLowerLimit,
			CreatedAtUpperLimit: params.Limitations.CreatedAtUpperLimit,
		})

	c.Handler.SetMaxQueryLimit(params.Limitations.MaxQueryLimit)
	c.Handler.SetDefaultQueryLimit(params.Limitations.DefaultQueryLimit)

	return nil
}
