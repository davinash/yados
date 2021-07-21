package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// GetListOfPeers Get the list of peers in the cluster
func (server *server) GetListOfPeers(ctx context.Context, request *pb.ListOfPeersRequest) (*pb.ListOfPeersReply, error) {
	servers := make([]*pb.Member, 0)
	for _, peer := range server.peers {
		servers = append(servers, &pb.Member{
			Name:    peer.Name,
			Address: peer.Address,
			Port:    peer.Port,
		})
	}
	// Add self
	servers = append(servers, server.self)

	reply := &pb.ListOfPeersReply{
		Member: servers,
	}

	return reply, nil
}
