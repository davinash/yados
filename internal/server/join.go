package server

import (
	"context"
	"fmt"
	pb "github.com/davinash/yados/internal/proto/gen"
)

//AddNewMemberInCluster adds a new member in cluster
func (server *YadosServer) AddNewMemberInCluster(ctx context.Context, newPeer *pb.NewMemberRequest) (*pb.NewMemberReply, error) {
	server.logger.Info("Adding new member in the cluster")

	if peer, ok := server.peers[newPeer.Member.Name]; ok {
		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
	}

	// Add new member in self
	server.peers[newPeer.Member.Name] = &pb.Member{
		Port:    newPeer.Member.Port,
		Address: newPeer.Member.Address,
		Name:    newPeer.Member.Name,
	}

	return &pb.NewMemberReply{Member: &pb.Member{
		Name:    server.self.Name,
		Address: server.self.Address,
		Port:    server.self.Port,
	}}, nil
}
