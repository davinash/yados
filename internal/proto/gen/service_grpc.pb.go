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

// YadosServiceClient is the client API for YadosService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type YadosServiceClient interface {
	RequestVotes(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error)
	AppendEntries(ctx context.Context, in *AppendEntryRequest, opts ...grpc.CallOption) (*AppendEntryReply, error)
	AddMember(ctx context.Context, in *NewPeerRequest, opts ...grpc.CallOption) (*NewPeerReply, error)
}

type yadosServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewYadosServiceClient(cc grpc.ClientConnInterface) YadosServiceClient {
	return &yadosServiceClient{cc}
}

func (c *yadosServiceClient) RequestVotes(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteReply, error) {
	out := new(VoteReply)
	err := c.cc.Invoke(ctx, "/YadosService/RequestVotes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) AppendEntries(ctx context.Context, in *AppendEntryRequest, opts ...grpc.CallOption) (*AppendEntryReply, error) {
	out := new(AppendEntryReply)
	err := c.cc.Invoke(ctx, "/YadosService/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) AddMember(ctx context.Context, in *NewPeerRequest, opts ...grpc.CallOption) (*NewPeerReply, error) {
	out := new(NewPeerReply)
	err := c.cc.Invoke(ctx, "/YadosService/AddMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// YadosServiceServer is the server API for YadosService service.
// All implementations must embed UnimplementedYadosServiceServer
// for forward compatibility
type YadosServiceServer interface {
	RequestVotes(context.Context, *VoteRequest) (*VoteReply, error)
	AppendEntries(context.Context, *AppendEntryRequest) (*AppendEntryReply, error)
	AddMember(context.Context, *NewPeerRequest) (*NewPeerReply, error)
	mustEmbedUnimplementedYadosServiceServer()
}

// UnimplementedYadosServiceServer must be embedded to have forward compatible implementations.
type UnimplementedYadosServiceServer struct {
}

func (UnimplementedYadosServiceServer) RequestVotes(context.Context, *VoteRequest) (*VoteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVotes not implemented")
}
func (UnimplementedYadosServiceServer) AppendEntries(context.Context, *AppendEntryRequest) (*AppendEntryReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedYadosServiceServer) AddMember(context.Context, *NewPeerRequest) (*NewPeerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMember not implemented")
}
func (UnimplementedYadosServiceServer) mustEmbedUnimplementedYadosServiceServer() {}

// UnsafeYadosServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to YadosServiceServer will
// result in compilation errors.
type UnsafeYadosServiceServer interface {
	mustEmbedUnimplementedYadosServiceServer()
}

func RegisterYadosServiceServer(s grpc.ServiceRegistrar, srv YadosServiceServer) {
	s.RegisterService(&YadosService_ServiceDesc, srv)
}

func _YadosService_RequestVotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).RequestVotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/RequestVotes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).RequestVotes(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).AppendEntries(ctx, req.(*AppendEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_AddMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewPeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).AddMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/AddMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).AddMember(ctx, req.(*NewPeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// YadosService_ServiceDesc is the grpc.ServiceDesc for YadosService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var YadosService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "YadosService",
	HandlerType: (*YadosServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestVotes",
			Handler:    _YadosService_RequestVotes_Handler,
		},
		{
			MethodName: "AppendEntries",
			Handler:    _YadosService_AppendEntries_Handler,
		},
		{
			MethodName: "AddMember",
			Handler:    _YadosService_AddMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
