package tests

import (
	"fmt"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
)

func createStore(srv server.Server, storeName string) (*pb.StoreCreateReply, error) {
	return server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: storeName,
	}, srv.Address(), srv.Port())
}

func performPut(srv server.Server, numOfPuts int, storeName string, prefix string) error {
	for i := 0; i < numOfPuts; i++ {
		err := server.ExecuteCmdPut(&server.PutArgs{
			Key:       fmt.Sprintf("%s-Key-%d", prefix, i),
			Value:     fmt.Sprintf("%s-Value-%d", prefix, i),
			StoreName: storeName,
		}, srv.Address(), srv.Port())
		if err != nil {
			return err
		}
	}
	return nil
}

//
//func (suite *YadosTestSuite) TestRandomServerDown() {
//	storeName := "TestRandomServerDown"
//	//WaitForLeaderElection(suite.cluster)
//
//	if _, err := createStore(suite.cluster.members[0], storeName); err != nil {
//		suite.T().Fatal(err)
//	}
//	if err := performPut(suite.cluster.members[0], 10, storeName, "BeforeRestart"); err != nil {
//		suite.T().Fatal(err)
//	}
//	// Let's stop Server
//	err := suite.cluster.members[0].Stop()
//	if err != nil {
//		suite.T().Fatal(err)
//	}
//	//WaitForLeaderElection(suite.cluster)
//	if err := performPut(suite.cluster.members[1], 10, storeName, "AfterOneNodeDown"); err != nil {
//		suite.T().Fatal(err)
//	}
//	//suite.cluster.members[0].Serve( /*[]*pb.Peer{suite.cluster.members[1].Self(),
//	//suite.cluster.members[2].Self()}*/)
//}
