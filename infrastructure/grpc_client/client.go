package grpcclient

import (
	"context"

	mpb "github.com/dezh-tech/immortal/infrastructure/grpc_client/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	RegistryService   mpb.ServiceRegistryClient
	ParametersService mpb.ParametersClient
	conn              *grpc.ClientConn
}

func New(endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		RegistryService:   mpb.NewServiceRegistryClient(conn),
		ParametersService: mpb.NewParametersClient(conn),
		conn:              conn,
	}, nil
}

func (c *Client) RegisterService(ctx context.Context,
	port, region string, hb uint32,
) (*mpb.RegisterServiceResponse, error) {
	return c.RegistryService.RegisterService(ctx, &mpb.RegisterServiceRequest{
		Type:                   mpb.ServiceTypeEnum_RELAY,
		Port:                   port,
		HeartbeatDurationInSec: hb,
		Region:                 region,
	})
}

func (c *Client) GetParameters(ctx context.Context, id string) (*mpb.GetParametersResponse, error) {
	md := metadata.New(map[string]string{"x-identifier": id})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return c.ParametersService.GetParameters(ctx, &mpb.GetParametersRequest{})
}
