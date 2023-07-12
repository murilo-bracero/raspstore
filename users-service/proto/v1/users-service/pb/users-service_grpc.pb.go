// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// UserConfigServiceClient is the client API for UserConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserConfigServiceClient interface {
	GetUserConfiguration(ctx context.Context, in *GetUserConfigurationRequest, opts ...grpc.CallOption) (*UserConfiguration, error)
}

type userConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserConfigServiceClient(cc grpc.ClientConnInterface) UserConfigServiceClient {
	return &userConfigServiceClient{cc}
}

func (c *userConfigServiceClient) GetUserConfiguration(ctx context.Context, in *GetUserConfigurationRequest, opts ...grpc.CallOption) (*UserConfiguration, error) {
	out := new(UserConfiguration)
	err := c.cc.Invoke(ctx, "/pb.UserConfigService/getUserConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserConfigServiceServer is the server API for UserConfigService service.
// All implementations must embed UnimplementedUserConfigServiceServer
// for forward compatibility
type UserConfigServiceServer interface {
	GetUserConfiguration(context.Context, *GetUserConfigurationRequest) (*UserConfiguration, error)
	mustEmbedUnimplementedUserConfigServiceServer()
}

// UnimplementedUserConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserConfigServiceServer struct {
}

func (UnimplementedUserConfigServiceServer) GetUserConfiguration(context.Context, *GetUserConfigurationRequest) (*UserConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserConfiguration not implemented")
}
func (UnimplementedUserConfigServiceServer) mustEmbedUnimplementedUserConfigServiceServer() {}

// UnsafeUserConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserConfigServiceServer will
// result in compilation errors.
type UnsafeUserConfigServiceServer interface {
	mustEmbedUnimplementedUserConfigServiceServer()
}

func RegisterUserConfigServiceServer(s grpc.ServiceRegistrar, srv UserConfigServiceServer) {
	s.RegisterService(&UserConfigService_ServiceDesc, srv)
}

func _UserConfigService_GetUserConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserConfigServiceServer).GetUserConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.UserConfigService/getUserConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserConfigServiceServer).GetUserConfiguration(ctx, req.(*GetUserConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserConfigService_ServiceDesc is the grpc.ServiceDesc for UserConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserConfigService",
	HandlerType: (*UserConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getUserConfiguration",
			Handler:    _UserConfigService_GetUserConfiguration_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/users-service.proto",
}
