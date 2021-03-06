// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package userproto

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

// UserSvcClient is the client API for UserSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserSvcClient interface {
	GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserRes, error)
}

type userSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewUserSvcClient(cc grpc.ClientConnInterface) UserSvcClient {
	return &userSvcClient{cc}
}

func (c *userSvcClient) GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserRes, error) {
	out := new(GetUserRes)
	err := c.cc.Invoke(ctx, "/userproto.UserSvc/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserSvcServer is the server API for UserSvc service.
// All implementations must embed UnimplementedUserSvcServer
// for forward compatibility
type UserSvcServer interface {
	GetUser(context.Context, *GetUserReq) (*GetUserRes, error)
	mustEmbedUnimplementedUserSvcServer()
}

// UnimplementedUserSvcServer must be embedded to have forward compatible implementations.
type UnimplementedUserSvcServer struct {
}

func (UnimplementedUserSvcServer) GetUser(context.Context, *GetUserReq) (*GetUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserSvcServer) mustEmbedUnimplementedUserSvcServer() {}

// UnsafeUserSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserSvcServer will
// result in compilation errors.
type UnsafeUserSvcServer interface {
	mustEmbedUnimplementedUserSvcServer()
}

func RegisterUserSvcServer(s grpc.ServiceRegistrar, srv UserSvcServer) {
	s.RegisterService(&UserSvc_ServiceDesc, srv)
}

func _UserSvc_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSvcServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userproto.UserSvc/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSvcServer).GetUser(ctx, req.(*GetUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UserSvc_ServiceDesc is the grpc.ServiceDesc for UserSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "userproto.UserSvc",
	HandlerType: (*UserSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _UserSvc_GetUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "userproto/user.proto",
}
