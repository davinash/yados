package tests

import (
	"fmt"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestCmdClusterStatus() {
	status, err := server.ExecuteCmdStatus(suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	if len(status.PeerStatus) != 3 {
		suite.T().Fatalf("expected reply containing 3 peers, actual = %d", len(status.PeerStatus))
	}
}

func (suite *YadosTestSuite) TestCmdCreateStoreKV() {

	storeName := "TestCmdCreateStoreKV"
	resp, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: storeName,
		Type: "memory",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	fmt.Println(resp.Msg)

	// Store already exists
	resp, err = server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: storeName,
		Type: "memory",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err == nil {
		suite.T().Fatal(err)
	}
	fmt.Println(err)

}

//
//func (suite *YadosTestSuite) TestCmdCreateStoreKVParallel() {
//	storeName := "TestCmdCreateStoreKVParallel"
//	wg := sync.WaitGroup{}
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func(wg *sync.WaitGroup, index int) {
//			defer wg.Done()
//			resp, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
//				Name: fmt.Sprintf("%s-%d", storeName, index),
//				Type: "memory",
//			}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
//			if err != nil {
//				suite.T().Fatal(err)
//			}
//			fmt.Println(resp.Msg)
//		}(&wg, i)
//	}
//	wg.Wait()
//}

func (suite *YadosTestSuite) TestCmdListStores() {
}
func (suite *YadosTestSuite) TestCmdPut() {
}
func (suite *YadosTestSuite) TestCmdGet() {
}
func (suite *YadosTestSuite) TestCmdDeleteStore() {
}
func (suite *YadosTestSuite) TestCmdExecuteQuery() {
}
func (suite *YadosTestSuite) TestCmdQuery() {
}
