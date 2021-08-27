package tests

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/davinash/yados/internal/events"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/davinash/yados/internal/server"
)

func printEntry(entry *pb.LogEntry, prefix string, t *testing.T) {
	t.Logf("[%s] Term= %v, Index= %v, Id= %v Type= %v", prefix, entry.Term, entry.Index, entry.Id, entry.CmdType)
	switch entry.CmdType {
	case pb.CommandType_CreateStore:
		var req pb.StoreCreateRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return
		}
		marshal, err := json.Marshal(&req)
		if err != nil {
			return
		}
		t.Logf("[%s] Command -> %s", prefix, string(marshal))

	case pb.CommandType_Put:
		var req pb.PutRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return
		}
		marshal, err := json.Marshal(&req)
		if err != nil {
			return
		}
		t.Logf("[%s] Command -> %s", prefix, string(marshal))

	case pb.CommandType_DeleteStore:
		var req pb.StoreDeleteRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return
		}
		marshal, err := json.Marshal(&req)
		if err != nil {
			return
		}
		t.Logf("[%s] Command -> %s", prefix, string(marshal))
	}
}

func (suite *YadosTestSuite) TestPLogAppend() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestPLogAppend"

	numOfPuts := 10

	for _, s := range suite.cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(events.EntryPersistEvents)
	}
	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(events.EntryPersistEvents)
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

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	for i := 0; i < numOfPuts; i++ {
		err := server.ExecuteCmdPut(&server.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: storeName,
		})
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	// Wait for replication to happen
	wg.Wait()

	for _, member := range suite.cluster.members {
		iterator, err := member.PLog().Iterator()
		if err != nil {
			suite.T().Fatal(err)
		}
		entry, err1 := iterator.Next()
		if err1 != nil {
			suite.T().Fatal(err1)
		}
		count := 0
		for entry != nil {
			printEntry(entry, member.Name(), suite.T())
			count++
			entry, err = iterator.Next()
			if err != nil {
				suite.T().Fatal(err)
			}
		}
		if count != numOfPuts+1 {
			suite.T().Fatalf("[%s] Expected entries = %d Actual = %d", member.Name(), numOfPuts, count)
		}

		err = iterator.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}
}

func (suite *YadosTestSuite) TestPLogAppendVerifyEntries() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestPLogAppendVerifyEntries"

	numOfPuts := 10

	for _, s := range suite.cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(events.EntryPersistEvents)
	}
	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(events.EntryPersistEvents)
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

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	for i := 0; i < numOfPuts; i++ {
		err := server.ExecuteCmdPut(&server.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: storeName,
		})
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	// Wait for replication to happen
	wg.Wait()

	for _, member := range suite.cluster.members {
		iterator, err := member.PLog().Iterator()
		if err != nil {
			suite.T().Fatal(err)
		}
		entry, err1 := iterator.Next()
		if err1 != nil {
			suite.T().Fatal(err1)
		}
		entry, err1 = iterator.Next()
		if err1 != nil {
			suite.T().Fatal(err1)
		}
		keyIdx := 0
		for entry != nil {
			var command pb.PutRequest
			err = anypb.UnmarshalTo(entry.Command, &command, proto.UnmarshalOptions{})
			if err != nil {
				suite.T().Fatal(err)
			}

			if command.Key != fmt.Sprintf("Key-%d", keyIdx) {
				suite.T().Fatalf("[%s] Expected Key = %s Actual = %s",
					member.Name(), fmt.Sprintf("Key-%d", keyIdx), command.Key)
			}
			if command.Value != fmt.Sprintf("Value-%d", keyIdx) {
				suite.T().Fatalf("[%s] Expected Values = %s Actual = %s",
					member.Name(), fmt.Sprintf("Value-%d", keyIdx), command.Value)
			}
			keyIdx++
			entry, err = iterator.Next()
			if err != nil {
				suite.T().Fatal(err)
			}
		}

		err = iterator.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}
}
