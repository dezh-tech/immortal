package client

import (
	"context"

	kraken "github.com/dezh-tech/immortal/client/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	RegistryService kraken.RegistryClient
	conn            *grpc.ClientConn
}

func NewClient(endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		RegistryService: kraken.NewRegistryClient(conn),
		conn:            conn,
	}, nil
}

func (c *Client) Register(ctx context.Context, url string, hb uint32) (*kraken.RegisterServiceResponse, error) {
	return c.RegistryService.RegisterService(ctx, &kraken.RegisterServiceRequest{
		Type:                   kraken.ServiceTypeEnum_RELAY,
		Url:                    url,
		HeartbeatDurationInSec: hb,
	})
}
