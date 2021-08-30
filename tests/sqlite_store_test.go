package tests

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreateSqlite() {
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
					if count == 2 {
						return
					}
				default:

				}
			}
		}(member)
	}

	storeName := "TestStoreCreateSqlite"

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
		Type:    "sqlite",
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	reply, err := server.ExecuteDDLQuery(&server.QueryArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		StoreName: storeName,
		SQLStr:    "create table employee(empid integer,name varchar(20),title varchar(10))",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	wg.Wait()

	if reply.Error != "" {
		suite.T().Fatalf("error message not expected, error = %s", reply.Error)
	}
	if reply.RowsAffected != 0 {
		suite.T().Fatalf("Expected 1, Actual %v", reply.RowsAffected)
	}

}
