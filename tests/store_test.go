package tests

import (
	"context"
	"log"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"google.golang.org/grpc"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	suite.WaitForLeaderElection()

	peerConn, rpcClient, err := server.GetPeerConn(suite.cluster.members[0].Address(),
		suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Error(err)
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	arg := &pb.StoreCreateRequest{Name: "TestStoreCreate"}
	_, err = rpcClient.CreateStore(context.Background(), arg)
	if err != nil {
		suite.T().Error(err)
	}
	time.Sleep(10 * time.Second)
}
