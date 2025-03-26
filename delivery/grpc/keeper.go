package grpc

import (
	"github.com/dezh-tech/immortal/delivery/websocket/configs"
	mpb "github.com/dezh-tech/immortal/infrastructure/grpc_client/gen"
	"github.com/dezh-tech/immortal/pkg/utils"
	"github.com/dezh-tech/immortal/repository"
)

type ParametersKeeper struct {
	Handler         *repository.Config
	WebsocketServer *configs.Config
}

func (keeper *ParametersKeeper) LoadParameters(params *mpb.GetParametersResponse) error {
	url, err := utils.ParseURL(params.Url)
	if err != nil {
		return err
	}

	keeper.WebsocketServer.SetURL(url)

	keeper.WebsocketServer.SetLimitation(
		&configs.Limitation{
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

	keeper.Handler.SetMaxQueryLimit(params.Limitations.MaxQueryLimit)
	keeper.Handler.SetDefaultQueryLimit(params.Limitations.DefaultQueryLimit)

	return nil
}
