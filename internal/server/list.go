package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// GetListOfPeers Get the list of peers in the cluster
func (server *YadosServer) GetListOfPeers(ctx context.Context, request *pb.ListOfPeersRequest) (*pb.ListOfPeersReply, error) {
	reply := &pb.ListOfPeersReply{
		Member: make([]*pb.Member, 0),
	}
	for _, p := range server.peers {
		reply.Member = append(reply.Member, &pb.Member{
			Name:    p.Name,
			Address: p.Address,
			Port:    p.Port,
		})
	}
	// Add Self
	reply.Member = append(reply.Member, server.self)
	return reply, nil
}
