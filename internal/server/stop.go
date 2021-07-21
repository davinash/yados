package server

import (
	"context"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// RemoveServerFromCluster Remove server from the cluster
func RemoveServerFromCluster(server Server, address string, port int32) error {
	conn, p, err := GetPeerConn(address, port)
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	_, err = p.RemoveServer(context.Background(), &pb.RemoveServerRequest{Member: &pb.Member{
		Name:    server.Self().Name,
		Address: server.Self().Address,
		Port:    server.Self().Port,
	}})
	if err != nil {
		return err
	}
	return nil
}

// StopServer Stops the server
func (server *server) StopServer(ctx context.Context, request *pb.StopServerRequest) (*pb.StopServerReply, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for _, s := range server.peers {
		err := RemoveServerFromCluster(server, s.Address, s.Port)
		if err != nil {
			server.logger.Warningf("Error removing server from cluster, Error = %v", err)
		}
	}
	if !server.isTestMode {
		server.osSignalCh <- syscall.SIGINT
	}
	return &pb.StopServerReply{}, nil
}
