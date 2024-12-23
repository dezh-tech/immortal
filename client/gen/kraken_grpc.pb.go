// Define the syntax version for Protocol Buffers

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: kraken.proto

// Define the package for the services and messages

package client

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
	ServiceRegistry_RegisterService_FullMethodName = "/kraken.ServiceRegistry/registerService"
)

// ServiceRegistryClient is the client API for ServiceRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceRegistryClient interface {
	// Registers a service with the registry
	RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error)
}

type serviceRegistryClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceRegistryClient(cc grpc.ClientConnInterface) ServiceRegistryClient {
	return &serviceRegistryClient{cc}
}

func (c *serviceRegistryClient) RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error) {
	out := new(RegisterServiceResponse)
	err := c.cc.Invoke(ctx, ServiceRegistry_RegisterService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceRegistryServer is the server API for ServiceRegistry service.
// All implementations should embed UnimplementedServiceRegistryServer
// for forward compatibility
type ServiceRegistryServer interface {
	// Registers a service with the registry
	RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error)
}

// UnimplementedServiceRegistryServer should be embedded to have forward compatible implementations.
type UnimplementedServiceRegistryServer struct {
}

func (UnimplementedServiceRegistryServer) RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterService not implemented")
}

// UnsafeServiceRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceRegistryServer will
// result in compilation errors.
type UnsafeServiceRegistryServer interface {
	mustEmbedUnimplementedServiceRegistryServer()
}

func RegisterServiceRegistryServer(s grpc.ServiceRegistrar, srv ServiceRegistryServer) {
	s.RegisterService(&ServiceRegistry_ServiceDesc, srv)
}

func _ServiceRegistry_RegisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceRegistryServer).RegisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServiceRegistry_RegisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceRegistryServer).RegisterService(ctx, req.(*RegisterServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServiceRegistry_ServiceDesc is the grpc.ServiceDesc for ServiceRegistry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServiceRegistry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kraken.ServiceRegistry",
	HandlerType: (*ServiceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "registerService",
			Handler:    _ServiceRegistry_RegisterService_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kraken.proto",
}

const (
	Config_GetConfig_FullMethodName = "/kraken.Config/getConfig"
)

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigClient interface {
	// Retrieves the current configuration
	GetConfig(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) GetConfig(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, Config_GetConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
// All implementations should embed UnimplementedConfigServer
// for forward compatibility
type ConfigServer interface {
	// Retrieves the current configuration
	GetConfig(context.Context, *EmptyRequest) (*GetConfigResponse, error)
}

// UnimplementedConfigServer should be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (UnimplementedConfigServer) GetConfig(context.Context, *EmptyRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}

// UnsafeConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServer will
// result in compilation errors.
type UnsafeConfigServer interface {
	mustEmbedUnimplementedConfigServer()
}

func RegisterConfigServer(s grpc.ServiceRegistrar, srv ConfigServer) {
	s.RegisterService(&Config_ServiceDesc, srv)
}

func _Config_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Config_GetConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetConfig(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Config_ServiceDesc is the grpc.ServiceDesc for Config service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Config_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kraken.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getConfig",
			Handler:    _Config_GetConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kraken.proto",
}

const (
	Log_AddLog_FullMethodName = "/kraken.Log/addLog"
)

// LogClient is the client API for Log service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogClient interface {
	// Adds a log entry
	AddLog(ctx context.Context, in *AddLogRequest, opts ...grpc.CallOption) (*AddLogResponse, error)
}

type logClient struct {
	cc grpc.ClientConnInterface
}

func NewLogClient(cc grpc.ClientConnInterface) LogClient {
	return &logClient{cc}
}

func (c *logClient) AddLog(ctx context.Context, in *AddLogRequest, opts ...grpc.CallOption) (*AddLogResponse, error) {
	out := new(AddLogResponse)
	err := c.cc.Invoke(ctx, Log_AddLog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServer is the server API for Log service.
// All implementations should embed UnimplementedLogServer
// for forward compatibility
type LogServer interface {
	// Adds a log entry
	AddLog(context.Context, *AddLogRequest) (*AddLogResponse, error)
}

// UnimplementedLogServer should be embedded to have forward compatible implementations.
type UnimplementedLogServer struct {
}

func (UnimplementedLogServer) AddLog(context.Context, *AddLogRequest) (*AddLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLog not implemented")
}

// UnsafeLogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogServer will
// result in compilation errors.
type UnsafeLogServer interface {
	mustEmbedUnimplementedLogServer()
}

func RegisterLogServer(s grpc.ServiceRegistrar, srv LogServer) {
	s.RegisterService(&Log_ServiceDesc, srv)
}

func _Log_AddLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServer).AddLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Log_AddLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServer).AddLog(ctx, req.(*AddLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Log_ServiceDesc is the grpc.ServiceDesc for Log service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Log_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kraken.Log",
	HandlerType: (*LogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "addLog",
			Handler:    _Log_AddLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kraken.proto",
}
