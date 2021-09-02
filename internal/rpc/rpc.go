package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//ErrorPortNotOpen returned when the Port is not open
var ErrorPortNotOpen = errors.New("port is not open to listen")

//ErrorUnknownMethod when rpc is called with unknown method
var ErrorUnknownMethod = errors.New("unknown method ")

//Server interface for rpc server
type Server interface {
	Start() error
	Stop() error
	Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error)
	GrpcServer() *grpc.Server
}

type rpcServer struct {
	grpcServer *grpc.Server
	logger     *logrus.Logger
	srvName    string
	address    string
	port       int32
}

//NewRPCServer creates a new instance of rpc server
func NewRPCServer(srvName string, address string, port int32, logger *logrus.Logger) Server {
	rpc := &rpcServer{
		grpcServer: grpc.NewServer(),
		logger:     logger,
		srvName:    srvName,
		address:    address,
		port:       port,
	}
	return rpc
}

func (rpc *rpcServer) GrpcServer() *grpc.Server {
	return rpc.grpcServer
}

func (rpc *rpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", rpc.address, rpc.port))
	if err != nil {
		return ErrorPortNotOpen
	}

	go func() {
		err := rpc.grpcServer.Serve(lis)
		if err != nil {
			rpc.logger.Errorf("failed to start the grpc server, Error = %v", err)
			return
		}
	}()
	return nil
}

func (rpc *rpcServer) Stop() error {
	rpc.logger.Debug("Stopping RPC server")
	rpc.grpcServer.Stop()
	rpc.logger.Debug("Stopped RPC server")
	return nil
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
			rpc.logger.Warnf("failed to close the connection, error = %v", err)
		}
	}(peerConn)

	switch serviceMethod {
	case "RPC.RequestVote":
		request := args.(*pb.VoteRequest)
		rpc.logger.Debugf("[%s] Type = RequestVotes %s -----> %s  "+
			"Request : {Term = %v Candidate Name = %v }", request.Id, rpc.srvName, peer.Name,
			request.Term, request.CandidateName)

		reply, err := rpcClient.RequestVotes(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "server.AddNewMember":
		request := args.(*pb.NewPeerRequest)
		rpc.logger.Debugf("[%s] Type = AddNewMember %s -----> %s  ",
			request.Id, rpc.srvName, peer.Name)

		reply, err := rpcClient.AddMember(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "RPC.AppendEntries":
		request := args.(*pb.AppendEntryRequest)

		rpc.logger.Debugf("[%s] AppendEntries ( -> %s ) : nextIndex = %d Term = %v; LeaderName = "+
			"%v; PrevLogTerm = %v; PrevLogIndex = %v; LeaderCommit = %v ", request.Id,
			peer.Name, request.NextIndex, request.Term, request.Leader.Name, request.PrevLogTerm,
			request.PrevLogIndex, request.LeaderCommit)

		reply, err := rpcClient.AppendEntries(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "RPC.RemoveSelf":
		request := args.(*pb.RemovePeerRequest)
		request.Id = uuid.New().String()
		rpc.logger.Debugf("[%s] RemovePeerRequest ", request.Id)
		reply, err := rpcClient.RemovePeer(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	case "RPC.PeerStatus":
		request := args.(*pb.StatusRequest)
		request.Id = uuid.New().String()
		rpc.logger.Debugf("[%s] PeerStatus ", request.Id)
		reply, err := rpcClient.PeerStatus(context.Background(), request)
		if err != nil {
			return nil, err
		}
		return reply, nil
	default:
		return nil, ErrorUnknownMethod
	}
}
