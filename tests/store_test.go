package tests

import (
	"fmt"
	"sync"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	WaitForLeaderElection(suite.cluster)

	for _, s := range suite.cluster.members {
		s.EventHandler().Subscribe(server.CommitEntryEvents)
	}

	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(server.CommitEntryEvents)
		}
	}()

	wg := sync.WaitGroup{}
	for _, member := range suite.cluster.members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			<-s.EventHandler().CommitEntryEvent()
		}(member)
	}

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

	for _, s := range suite.cluster.members {
		s.EventHandler().Subscribe(server.CommitEntryEvents)
	}

	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(server.CommitEntryEvents)
		}
	}()

	for i := 0; i < 5; i++ {
		wg := sync.WaitGroup{}
		for _, member := range suite.cluster.members {
			wg.Add(1)
			go func(s server.Server) {
				defer wg.Done()
				<-s.EventHandler().CommitEntryEvent()
			}(member)
		}

		err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
			Address: suite.cluster.members[0].Address(),
			Port:    suite.cluster.members[0].Port(),
			Name:    fmt.Sprintf("TestStoreList-%d", i),
		})
		if err != nil {
			suite.T().Fatal(err)
		}
		wg.Wait()
	}

	storeList, err := server.ExecuteCmdListStore(&server.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	if len(storeList.Name) != 5 {
		suite.T().Fatal(err)
	}
}
