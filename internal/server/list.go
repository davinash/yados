package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//GetListOfPeersEx helper function to get list of peers from this server
func (server *YadosServer) GetListOfPeersEx() []*pb.Member {
	servers := make([]*pb.Member, 0)
	for _, peer := range server.peers {
		servers = append(servers, &pb.Member{
			Name:    peer.Name,
			Address: peer.Address,
			Port:    peer.Port,
		})
	}
	// Add Self
	servers = append(servers, server.self)
	return servers
}

// GetListOfPeers Get the list of peers in the cluster
func (server *YadosServer) GetListOfPeers(ctx context.Context, request *pb.ListOfPeersRequest) (*pb.ListOfPeersReply, error) {
	reply := &pb.ListOfPeersReply{
		Member: server.GetListOfPeersEx(),
	}
	return reply, nil
}
