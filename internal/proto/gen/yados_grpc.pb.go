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
	AddNewMemberInCluster(ctx context.Context, in *NewMemberRequest, opts ...grpc.CallOption) (*NewMemberReply, error)
	GetListOfPeers(ctx context.Context, in *ListOfPeersRequest, opts ...grpc.CallOption) (*ListOfPeersReply, error)
	StopServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*StopServerReply, error)
	RemoveServer(ctx context.Context, in *RemoveServerRequest, opts ...grpc.CallOption) (*RemoveServerReply, error)
	UpdateHealthStatus(ctx context.Context, in *HealthStatusRequest, opts ...grpc.CallOption) (*HealthStatusReply, error)
	CreateStore(ctx context.Context, in *StoreCreateRequest, opts ...grpc.CallOption) (*StoreCreateReply, error)
	CreateStoreSecondary(ctx context.Context, in *StoreCreateRequest, opts ...grpc.CallOption) (*StoreCreateReply, error)
	RequestForVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteRequestReply, error)
}

type yadosServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewYadosServiceClient(cc grpc.ClientConnInterface) YadosServiceClient {
	return &yadosServiceClient{cc}
}

func (c *yadosServiceClient) AddNewMemberInCluster(ctx context.Context, in *NewMemberRequest, opts ...grpc.CallOption) (*NewMemberReply, error) {
	out := new(NewMemberReply)
	err := c.cc.Invoke(ctx, "/YadosService/AddNewMemberInCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) GetListOfPeers(ctx context.Context, in *ListOfPeersRequest, opts ...grpc.CallOption) (*ListOfPeersReply, error) {
	out := new(ListOfPeersReply)
	err := c.cc.Invoke(ctx, "/YadosService/GetListOfPeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) StopServer(ctx context.Context, in *StopServerRequest, opts ...grpc.CallOption) (*StopServerReply, error) {
	out := new(StopServerReply)
	err := c.cc.Invoke(ctx, "/YadosService/StopServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) RemoveServer(ctx context.Context, in *RemoveServerRequest, opts ...grpc.CallOption) (*RemoveServerReply, error) {
	out := new(RemoveServerReply)
	err := c.cc.Invoke(ctx, "/YadosService/RemoveServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) UpdateHealthStatus(ctx context.Context, in *HealthStatusRequest, opts ...grpc.CallOption) (*HealthStatusReply, error) {
	out := new(HealthStatusReply)
	err := c.cc.Invoke(ctx, "/YadosService/UpdateHealthStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) CreateStore(ctx context.Context, in *StoreCreateRequest, opts ...grpc.CallOption) (*StoreCreateReply, error) {
	out := new(StoreCreateReply)
	err := c.cc.Invoke(ctx, "/YadosService/CreateStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) CreateStoreSecondary(ctx context.Context, in *StoreCreateRequest, opts ...grpc.CallOption) (*StoreCreateReply, error) {
	out := new(StoreCreateReply)
	err := c.cc.Invoke(ctx, "/YadosService/CreateStoreSecondary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *yadosServiceClient) RequestForVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteRequestReply, error) {
	out := new(VoteRequestReply)
	err := c.cc.Invoke(ctx, "/YadosService/RequestForVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// YadosServiceServer is the server API for YadosService service.
// All implementations must embed UnimplementedYadosServiceServer
// for forward compatibility
type YadosServiceServer interface {
	AddNewMemberInCluster(context.Context, *NewMemberRequest) (*NewMemberReply, error)
	GetListOfPeers(context.Context, *ListOfPeersRequest) (*ListOfPeersReply, error)
	StopServer(context.Context, *StopServerRequest) (*StopServerReply, error)
	RemoveServer(context.Context, *RemoveServerRequest) (*RemoveServerReply, error)
	UpdateHealthStatus(context.Context, *HealthStatusRequest) (*HealthStatusReply, error)
	CreateStore(context.Context, *StoreCreateRequest) (*StoreCreateReply, error)
	CreateStoreSecondary(context.Context, *StoreCreateRequest) (*StoreCreateReply, error)
	RequestForVote(context.Context, *VoteRequest) (*VoteRequestReply, error)
	mustEmbedUnimplementedYadosServiceServer()
}

// UnimplementedYadosServiceServer must be embedded to have forward compatible implementations.
type UnimplementedYadosServiceServer struct {
}

func (UnimplementedYadosServiceServer) AddNewMemberInCluster(context.Context, *NewMemberRequest) (*NewMemberReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNewMemberInCluster not implemented")
}
func (UnimplementedYadosServiceServer) GetListOfPeers(context.Context, *ListOfPeersRequest) (*ListOfPeersReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListOfPeers not implemented")
}
func (UnimplementedYadosServiceServer) StopServer(context.Context, *StopServerRequest) (*StopServerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopServer not implemented")
}
func (UnimplementedYadosServiceServer) RemoveServer(context.Context, *RemoveServerRequest) (*RemoveServerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveServer not implemented")
}
func (UnimplementedYadosServiceServer) UpdateHealthStatus(context.Context, *HealthStatusRequest) (*HealthStatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHealthStatus not implemented")
}
func (UnimplementedYadosServiceServer) CreateStore(context.Context, *StoreCreateRequest) (*StoreCreateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStore not implemented")
}
func (UnimplementedYadosServiceServer) CreateStoreSecondary(context.Context, *StoreCreateRequest) (*StoreCreateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStoreSecondary not implemented")
}
func (UnimplementedYadosServiceServer) RequestForVote(context.Context, *VoteRequest) (*VoteRequestReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestForVote not implemented")
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

func _YadosService_AddNewMemberInCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).AddNewMemberInCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/AddNewMemberInCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).AddNewMemberInCluster(ctx, req.(*NewMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_GetListOfPeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOfPeersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).GetListOfPeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/GetListOfPeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).GetListOfPeers(ctx, req.(*ListOfPeersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_StopServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).StopServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/StopServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).StopServer(ctx, req.(*StopServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_RemoveServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).RemoveServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/RemoveServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).RemoveServer(ctx, req.(*RemoveServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_UpdateHealthStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).UpdateHealthStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/UpdateHealthStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).UpdateHealthStatus(ctx, req.(*HealthStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_CreateStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).CreateStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/CreateStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).CreateStore(ctx, req.(*StoreCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_CreateStoreSecondary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).CreateStoreSecondary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/CreateStoreSecondary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).CreateStoreSecondary(ctx, req.(*StoreCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _YadosService_RequestForVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YadosServiceServer).RequestForVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/YadosService/RequestForVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YadosServiceServer).RequestForVote(ctx, req.(*VoteRequest))
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
			MethodName: "AddNewMemberInCluster",
			Handler:    _YadosService_AddNewMemberInCluster_Handler,
		},
		{
			MethodName: "GetListOfPeers",
			Handler:    _YadosService_GetListOfPeers_Handler,
		},
		{
			MethodName: "StopServer",
			Handler:    _YadosService_StopServer_Handler,
		},
		{
			MethodName: "RemoveServer",
			Handler:    _YadosService_RemoveServer_Handler,
		},
		{
			MethodName: "UpdateHealthStatus",
			Handler:    _YadosService_UpdateHealthStatus_Handler,
		},
		{
			MethodName: "CreateStore",
			Handler:    _YadosService_CreateStore_Handler,
		},
		{
			MethodName: "CreateStoreSecondary",
			Handler:    _YadosService_CreateStoreSecondary_Handler,
		},
		{
			MethodName: "RequestForVote",
			Handler:    _YadosService_RequestForVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "yados.proto",
}
