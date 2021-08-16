package tests

import (
	"fmt"
	"sync"

	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestPLogAppend() {
	suite.WaitForLeaderElection()
	storeName := "TestPLogAppend"

	numOfPuts := 10

	for _, s := range suite.cluster.members {
		s.EventHandler().Subscribe(server.EntryPersistEvents)
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
	}

	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(server.EntryPersistEvents)
		}
	}()

	wg := sync.WaitGroup{}
	for _, member := range suite.cluster.members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			<-s.EventHandler().PersistEntryEvent()
		}(member)
	}

	err := store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		suite.T().Error(err)
	}

	for i := 0; i < numOfPuts; i++ {
		err := store.ExecutePutCommand(&store.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: storeName,
		})
		if err != nil {
			suite.T().Error(err)
		}
	}
	// Wait for replication to happen
	wg.Wait()

}
