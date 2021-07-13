package server

import (
	"context"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// RemoveServerFromCluster Remove server from the cluster
func (server *YadosServer) RemoveServerFromCluster(address string, port int32) error {
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
		Name:    server.self.Name,
		Address: server.self.Address,
		Port:    server.self.Port,
	}})
	if err != nil {
		return err
	}
	return nil
}

// StopServer Stops the server
func (server *YadosServer) StopServer(ctx context.Context, request *pb.StopServerRequest) (*pb.StopServerReply, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for _, s := range server.peers {
		err := server.RemoveServerFromCluster(s.Address, s.Port)
		if err != nil {
			server.logger.Warningf("Error removing server from cluster, Error = %v", err)
		}
	}
	if !server.isTestMode {
		server.OSSignalCh <- syscall.SIGINT
	}
	return &pb.StopServerReply{}, nil
}
