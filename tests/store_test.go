package tests

import (
	"fmt"
	"sync"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreate() {

	wg := sync.WaitGroup{}
	WaitForEvents(suite.cluster.members, &wg, 1)
	defer func() {
		StopWaitForEvents(suite.cluster.members)
	}()

	_, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: "TestStoreCreate",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	wg.Wait()
}

func (suite *YadosTestSuite) TestStoreList() {

	wg := sync.WaitGroup{}
	WaitForEvents(suite.cluster.members, &wg, 5)
	defer func() {
		StopWaitForEvents(suite.cluster.members)
	}()

	for i := 0; i < 5; i++ {
		_, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
			Name: fmt.Sprintf("TestStoreList-%d", i),
		}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	wg.Wait()

	storeList, err := server.ExecuteCmdListStore(suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}

	if len(storeList.Name) != 5 {
		suite.T().Fatal(err)
	}
}
