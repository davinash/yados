// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gen

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

// ControllerServiceClient is the client API for ControllerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControllerServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error)
	UnRegister(ctx context.Context, in *UnRegisterRequest, opts ...grpc.CallOption) (*UnRegisterReply, error)
	GetLeader(ctx context.Context, in *GetLeaderRequest, opts ...grpc.CallOption) (*GetLeaderReply, error)
}

type controllerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewControllerServiceClient(cc grpc.ClientConnInterface) ControllerServiceClient {
	return &controllerServiceClient{cc}
}

func (c *controllerServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterReply, error) {
	out := new(RegisterReply)
	err := c.cc.Invoke(ctx, "/ControllerService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerServiceClient) UnRegister(ctx context.Context, in *UnRegisterRequest, opts ...grpc.CallOption) (*UnRegisterReply, error) {
	out := new(UnRegisterReply)
	err := c.cc.Invoke(ctx, "/ControllerService/UnRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controllerServiceClient) GetLeader(ctx context.Context, in *GetLeaderRequest, opts ...grpc.CallOption) (*GetLeaderReply, error) {
	out := new(GetLeaderReply)
	err := c.cc.Invoke(ctx, "/ControllerService/GetLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControllerServiceServer is the server API for ControllerService service.
// All implementations must embed UnimplementedControllerServiceServer
// for forward compatibility
type ControllerServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterReply, error)
	UnRegister(context.Context, *UnRegisterRequest) (*UnRegisterReply, error)
	GetLeader(context.Context, *GetLeaderRequest) (*GetLeaderReply, error)
	mustEmbedUnimplementedControllerServiceServer()
}

// UnimplementedControllerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedControllerServiceServer struct {
}

func (UnimplementedControllerServiceServer) Register(context.Context, *RegisterRequest) (*RegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedControllerServiceServer) UnRegister(context.Context, *UnRegisterRequest) (*UnRegisterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnRegister not implemented")
}
func (UnimplementedControllerServiceServer) GetLeader(context.Context, *GetLeaderRequest) (*GetLeaderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeader not implemented")
}
func (UnimplementedControllerServiceServer) mustEmbedUnimplementedControllerServiceServer() {}

// UnsafeControllerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControllerServiceServer will
// result in compilation errors.
type UnsafeControllerServiceServer interface {
	mustEmbedUnimplementedControllerServiceServer()
}

func RegisterControllerServiceServer(s grpc.ServiceRegistrar, srv ControllerServiceServer) {
	s.RegisterService(&ControllerService_ServiceDesc, srv)
}

func _ControllerService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ControllerService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControllerService_UnRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServiceServer).UnRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ControllerService/UnRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServiceServer).UnRegister(ctx, req.(*UnRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControllerService_GetLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControllerServiceServer).GetLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ControllerService/GetLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControllerServiceServer).GetLeader(ctx, req.(*GetLeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ControllerService_ServiceDesc is the grpc.ServiceDesc for ControllerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ControllerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ControllerService",
	HandlerType: (*ControllerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _ControllerService_Register_Handler,
		},
		{
			MethodName: "UnRegister",
			Handler:    _ControllerService_UnRegister_Handler,
		},
		{
			MethodName: "GetLeader",
			Handler:    _ControllerService_GetLeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controller.proto",
}
