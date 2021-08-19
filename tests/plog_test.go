package tests

import (
	"fmt"
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestPLogAppend() {
	suite.WaitForLeaderElection()
	storeName := "TestPLogAppend"

	numOfPuts := 10

	for _, s := range suite.cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(server.EntryPersistEvents)
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

	for _, member := range suite.cluster.members {
		iterator, err := member.PLog().Iterator()
		if err != nil {
			suite.T().Error(err)
		}
		entry, err1 := iterator.Next()
		if err1 != nil {
			suite.T().Error(err1)
		}
		count := 0
		for entry != nil {
			count++
			entry, err = iterator.Next()
			if err != nil {
				suite.T().Error(err)
			}
		}
		if count != numOfPuts+1 {
			suite.T().Errorf("[%s] Expected entries = %d Actual = %d", member.Name(), numOfPuts, count)
		}

		err = iterator.Close()
		if err != nil {
			suite.T().Error(err)
		}
	}
}

func (suite *YadosTestSuite) TestPLogAppendVerifyEntries() {
	suite.WaitForLeaderElection()
	storeName := "TestPLogAppendVerifyEntries"

	numOfPuts := 10

	for _, s := range suite.cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(server.EntryPersistEvents)
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

	for _, member := range suite.cluster.members {
		iterator, err := member.PLog().Iterator()
		if err != nil {
			suite.T().Error(err)
		}
		entry, err1 := iterator.Next()
		if err1 != nil {
			suite.T().Error(err1)
		}
		entry, err1 = iterator.Next()
		if err1 != nil {
			suite.T().Error(err1)
		}
		keyIdx := 0
		for entry != nil {
			var command pb.PutRequest
			err = anypb.UnmarshalTo(entry.Command, &command, proto.UnmarshalOptions{})
			if err != nil {
				suite.T().Error(err)
			}

			if command.Key != fmt.Sprintf("Key-%d", keyIdx) {
				suite.T().Errorf("[%s] Expected Key = %s Actual = %s",
					member.Name(), fmt.Sprintf("Key-%d", keyIdx), command.Key)
			}
			if command.Value != fmt.Sprintf("Value-%d", keyIdx) {
				suite.T().Errorf("[%s] Expected Values = %s Actual = %s",
					member.Name(), fmt.Sprintf("Value-%d", keyIdx), command.Value)
			}
			keyIdx++
			entry, err = iterator.Next()
			if err != nil {
				suite.T().Error(err)
			}
		}

		err = iterator.Close()
		if err != nil {
			suite.T().Error(err)
		}
	}
}
