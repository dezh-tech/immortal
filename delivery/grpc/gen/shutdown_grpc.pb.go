// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: shutdown.proto

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
	ShutdownService_Shutdown_FullMethodName = "/relay.v1.ShutdownService/Shutdown"
)

// ShutdownServiceClient is the client API for ShutdownService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShutdownServiceClient interface {
	Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error)
}

type shutdownServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShutdownServiceClient(cc grpc.ClientConnInterface) ShutdownServiceClient {
	return &shutdownServiceClient{cc}
}

func (c *shutdownServiceClient) Shutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error) {
	out := new(ShutdownResponse)
	err := c.cc.Invoke(ctx, ShutdownService_Shutdown_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShutdownServiceServer is the server API for ShutdownService service.
// All implementations should embed UnimplementedShutdownServiceServer
// for forward compatibility
type ShutdownServiceServer interface {
	Shutdown(context.Context, *ShutdownRequest) (*ShutdownResponse, error)
}

// UnimplementedShutdownServiceServer should be embedded to have forward compatible implementations.
type UnimplementedShutdownServiceServer struct {
}

func (UnimplementedShutdownServiceServer) Shutdown(context.Context, *ShutdownRequest) (*ShutdownResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shutdown not implemented")
}

// UnsafeShutdownServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShutdownServiceServer will
// result in compilation errors.
type UnsafeShutdownServiceServer interface {
	mustEmbedUnimplementedShutdownServiceServer()
}

func RegisterShutdownServiceServer(s grpc.ServiceRegistrar, srv ShutdownServiceServer) {
	s.RegisterService(&ShutdownService_ServiceDesc, srv)
}

func _ShutdownService_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShutdownServiceServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShutdownService_Shutdown_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShutdownServiceServer).Shutdown(ctx, req.(*ShutdownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShutdownService_ServiceDesc is the grpc.ServiceDesc for ShutdownService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShutdownService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relay.v1.ShutdownService",
	HandlerType: (*ShutdownServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Shutdown",
			Handler:    _ShutdownService_Shutdown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shutdown.proto",
}
