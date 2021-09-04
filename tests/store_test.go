package tests

import (
	"fmt"
	"sync"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	WaitForLeaderElection(suite.cluster)

	wg := sync.WaitGroup{}
	WaitForEvents(suite.cluster.members, &wg, 1)
	defer func() {
		StopWaitForEvents(suite.cluster.members)
	}()

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestStoreCreate",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	wg.Wait()
}

func (suite *YadosTestSuite) TestStoreList() {
	WaitForLeaderElection(suite.cluster)

	wg := sync.WaitGroup{}
	WaitForEvents(suite.cluster.members, &wg, 5)
	defer func() {
		StopWaitForEvents(suite.cluster.members)
	}()

	for i := 0; i < 5; i++ {
		err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
			Address: suite.cluster.members[0].Address(),
			Port:    suite.cluster.members[0].Port(),
			Name:    fmt.Sprintf("TestStoreList-%d", i),
		})
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	wg.Wait()

	storeList, err := server.ExecuteCmdListStore(&server.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	if len(storeList.Stores) != 5 {
		suite.T().Fatal(err)
	}

	for _, s := range storeList.Stores {
		suite.T().Logf("Type => %v", s.StoreType)
	}
}
