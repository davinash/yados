package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//ErrorPortNotOpen returned when the port is not open
var ErrorPortNotOpen = errors.New("port is not open to listen")

//ErrorUnknownMethod when rpc is called with unknown method
var ErrorUnknownMethod = errors.New("unknown method ")

//RPCServer interface for rpc server
type RPCServer interface {
	Start() error
	Stop() error
	Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error)
	Server() Server
}

type rpcServer struct {
	server     Server
	grpcServer *grpc.Server
}

//NewRPCServer creates a new instance of rpc server
func NewRPCServer(srv Server) RPCServer {
	rpc := &rpcServer{
		grpcServer: grpc.NewServer(),
		server:     srv,
	}
	return rpc
}

func (rpc *rpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", rpc.server.Address(), rpc.server.Port()))
	if err != nil {
		return ErrorPortNotOpen
	}

	pb.RegisterYadosServiceServer(rpc.grpcServer, rpc.server)
	go func() {
		rpc.grpcServer.Serve(lis)
	}()
	return nil
}

func (rpc *rpcServer) Stop() error {
	rpc.grpcServer.Stop()
	return nil
}

func (rpc *rpcServer) Server() Server {
	return rpc.server
}

//GetPeerConn returns the connection and grpc client for the remote peer
func GetPeerConn(address string, port int32) (*grpc.ClientConn, pb.YadosServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	peer := pb.NewYadosServiceClient(conn)
	return conn, peer, nil
}

func (rpc *rpcServer) Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error) {
	peerConn, rpcClient, err := GetPeerConn(peer.Address, peer.Port)
	if err != nil {
		return nil, err
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			rpc.Server().Logger().Warnf("failed to close the connection, error = %v", err)
		}
	}(peerConn)

	switch serviceMethod {
	case "RPC.RequestVote":
		reply, err := rpcClient.RequestVotes(context.Background(), args.(*pb.VoteRequest))
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "server.AddNewMember":
		reply, err := rpcClient.AddMember(context.Background(), args.(*pb.Peer))
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "RPC.AppendEntries":
		reply, err := rpcClient.AppendEntries(context.Background(), args.(*pb.AppendEntryRequest))
		if err != nil {
			return nil, err
		}
		return reply, nil
	default:
		return nil, ErrorUnknownMethod
	}
}
