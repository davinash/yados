package server

import (
	"context"
	"fmt"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//RemoveServer Remove server from the cluster
func (server *YadosServer) RemoveServer(ctx context.Context, request *pb.RemoveServerRequest) (*pb.RemoveServerReply, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if _, ok := server.peers[request.Member.Name]; !ok {
		return &pb.RemoveServerReply{}, fmt.Errorf("server with Name %s does not exists in cluster", request.Member.Name)
	}
	delete(server.peers, request.Member.Name)

	return &pb.RemoveServerReply{}, nil
}
