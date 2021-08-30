package tests

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	WaitForLeaderElection(suite.cluster)

	for _, s := range suite.cluster.members {
		s.EventHandler().PersistEntryChan = make(chan *pb.LogEntry)
	}
	defer func() {
		for _, s := range suite.cluster.members {
			close(s.EventHandler().PersistEntryChan)
			s.EventHandler().PersistEntryChan = nil
		}
	}()

	wg := sync.WaitGroup{}
	for _, member := range suite.cluster.members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			count := 0
			for {
				select {
				case <-s.EventHandler().PersistEntryChan:
					count++
					if count == 1 {
						return
					}
				default:

				}
			}
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

	//for _, s := range suite.cluster.members {
	//	s.EventHandler().Subscribe(events.CommitEntryEvents)
	//}
	//
	//defer func() {
	//	for _, s := range suite.cluster.members {
	//		s.EventHandler().UnSubscribe(events.CommitEntryEvents)
	//	}
	//}()

	//for i := 0; i < 5; i++ {
	//	wg := sync.WaitGroup{}
	//	for _, member := range suite.cluster.members {
	//		wg.Add(1)
	//		go func(s server.Server) {
	//			defer wg.Done()
	//			<-s.EventHandler().CommitEntryEvent()
	//		}(member)
	//	}
	//
	//	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
	//		Address: suite.cluster.members[0].Address(),
	//		Port:    suite.cluster.members[0].Port(),
	//		Name:    fmt.Sprintf("TestStoreList-%d", i),
	//	})
	//	if err != nil {
	//		suite.T().Fatal(err)
	//	}
	//	wg.Wait()
	//}

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
