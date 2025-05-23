// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: health.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Health_Status_FullMethodName = "/relay.v1.Health/Status"
)

// HealthClient is the client API for Health service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthClient interface {
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type healthClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthClient(cc grpc.ClientConnInterface) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, Health_Status_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthServer is the server API for Health service.
// All implementations should embed UnimplementedHealthServer
// for forward compatibility
type HealthServer interface {
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
}

// UnimplementedHealthServer should be embedded to have forward compatible implementations.
type UnimplementedHealthServer struct {
}

func (UnimplementedHealthServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}

// UnsafeHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthServer will
// result in compilation errors.
type UnsafeHealthServer interface {
	mustEmbedUnimplementedHealthServer()
}

func RegisterHealthServer(s grpc.ServiceRegistrar, srv HealthServer) {
	s.RegisterService(&Health_ServiceDesc, srv)
}

func _Health_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Health_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Health_ServiceDesc is the grpc.ServiceDesc for Health service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Health_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relay.v1.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _Health_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "health.proto",
}
