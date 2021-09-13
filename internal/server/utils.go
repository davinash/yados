package server

import (
	"context"
	"github.com/davinash/yados/internal/rpc"
	"log"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//GetLeader returns the current leader in the system
func GetLeader(address string, port int32) (*pb.Peer, error) {
	peerConn, rpcClient, err := rpc.GetPeerConn(address, port)
	if err != nil {
		return nil, err
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

TryLeaderQueryAgain:
	leader, err := rpcClient.GetLeader(context.Background(), &pb.GetLeaderRequest{})
	if err != nil {
		return nil, err
	}
	if leader.Leader == nil {
		time.Sleep(10 * time.Millisecond)
		goto TryLeaderQueryAgain
	}
	return leader.Leader, nil
}
