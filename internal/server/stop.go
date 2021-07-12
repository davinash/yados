package server

import (
	"context"
	"syscall"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// StopServer Stops the server
func (server *YadosServer) StopServer(ctx context.Context, request *pb.StopServerRequest) (*pb.StopServerReply, error) {
	server.OSSignalCh <- syscall.SIGINT
	return &pb.StopServerReply{}, nil
}
