package server

import (
	"errors"
	"fmt"
	"net"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//ErrorPortNotOpen returned when the port is not open
var ErrorPortNotOpen = errors.New("port is not open to listen")

//RPCServer interface for rpc server
type RPCServer interface {
	Start() error
	Stop() error
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

func (srv *rpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.server.Address(), srv.server.Port()))
	if err != nil {
		return ErrorPortNotOpen
	}

	pb.RegisterYadosServiceServer(srv.grpcServer, srv.server)
	go func() {
		srv.grpcServer.Serve(lis)
	}()
	return nil
}

func (srv *rpcServer) Stop() error {
	srv.grpcServer.Stop()
	return nil
}
