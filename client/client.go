package client

import (
	"context"

	kraken "github.com/dezh-tech/immortal/client/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	RegistryService kraken.ServiceRegistryClient
	ConfigService   kraken.ConfigClient
	conn            *grpc.ClientConn
}

func NewClient(endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		RegistryService: kraken.NewServiceRegistryClient(conn),
		ConfigService:   kraken.NewConfigClient(conn),
		conn:            conn,
	}, nil
}

func (c *Client) RegisterService(ctx context.Context,
	port, region string, hb uint32,
) (*kraken.RegisterServiceResponse, error) {
	return c.RegistryService.RegisterService(ctx, &kraken.RegisterServiceRequest{
		Type:                   kraken.ServiceTypeEnum_RELAY,
		Port:                   port,
		HeartbeatDurationInSec: hb,
		Region:                 region,
	})
}

func (c *Client) GetConfig(ctx context.Context, id string) (*kraken.GetConfigResponse, error) {
	md := metadata.New(map[string]string{"x-identifier": id})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return c.ConfigService.GetConfig(ctx, &kraken.EmptyRequest{})
}
