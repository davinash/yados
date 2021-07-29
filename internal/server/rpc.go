package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"

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
		request := args.(*pb.VoteRequest)
		request.Id = uuid.New().String()
		rpc.Server().Logger().Debugf("[%s] Type = RequestVotes %s -----> %s  "+
			"Request : {Term = %v Candidate Name = %v }", request.Id, rpc.Server().Name(), peer.Name,
			request.Term, request.CandidateName)

		reply, err := rpcClient.RequestVotes(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "server.AddNewMember":
		request := args.(*pb.NewPeerRequest)
		request.Id = uuid.New().String()
		rpc.Server().Logger().Debugf("[%s] Type = AddNewMember %s -----> %s  ",
			request.Id, rpc.Server().Name(), peer.Name)

		reply, err := rpcClient.AddMember(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "RPC.AppendEntries":
		request := args.(*pb.AppendEntryRequest)
		request.Id = uuid.New().String()
		rpc.Server().Logger().Debugf("[%s] Type = AppendEntries %s -----> %s  "+
			"Request : {Term = %v Leader Name = %v }", request.Id, rpc.Server().Name(), peer.Name,
			request.Term, request.LeaderName)

		reply, err := rpcClient.AppendEntries(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	default:
		return nil, ErrorUnknownMethod
	}
}
