package server

import (
	"context"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//GetLeader returns the current leader in the system
func GetLeader(address string, port int32) (*pb.Peer, error) {
	peerConn, rpcClient, err := GetPeerConn(address, port)
	if err != nil {
		return nil, err
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	clusterStatus, err := rpcClient.ClusterStatus(context.Background(), &pb.ClusterStatusRequest{})
	if err != nil {
		return &pb.Peer{}, err
	}

	for _, peer := range clusterStatus.PeerStatus {
		if peer.IsLeader {
			return &pb.Peer{
				Address: peer.Server.Address,
				Port:    peer.Server.Port,
			}, nil
		}
	}

	return &pb.Peer{}, nil
}
