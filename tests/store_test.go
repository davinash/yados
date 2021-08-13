package tests

import (
	"fmt"
	"sync"

	"github.com/davinash/yados/cmd/cli/commands/store"
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	suite.WaitForLeaderElection()
	for _, s := range suite.cluster.members {
		s.SetTestArgs(&server.TestArgs{
			CommitEntryChan: make(chan *pb.LogEntry),
		})
	}

	wg := sync.WaitGroup{}
	for _, member := range suite.cluster.members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			<-s.TestArgs().CommitEntryChan
		}(member)
	}

	err := store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestStoreCreate",
	})
	if err != nil {
		suite.T().Error(err)
	}

	wg.Wait()

}

func (suite *YadosTestSuite) TestStoreList() {
	suite.WaitForLeaderElection()
	for _, s := range suite.cluster.members {
		s.SetTestArgs(&server.TestArgs{
			CommitEntryChan: make(chan *pb.LogEntry),
		})
	}

	for i := 0; i < 5; i++ {
		wg := sync.WaitGroup{}
		for _, member := range suite.cluster.members {
			wg.Add(1)
			go func(s server.Server) {
				defer wg.Done()
				<-s.TestArgs().CommitEntryChan
			}(member)
		}

		err := store.CreateCommandExecute(&store.CreateCommandArgs{
			Address: suite.cluster.members[0].Address(),
			Port:    suite.cluster.members[0].Port(),
			Name:    fmt.Sprintf("TestStoreList-%d", i),
		})
		if err != nil {
			suite.T().Error(err)
		}
		wg.Wait()
	}

	storeList, err := store.ExecuteStoreListCommand(&store.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		suite.T().Error(err)
	}
	if len(storeList.Name) != 5 {
		suite.T().Error(err)
	}
}
