package tests

import (
	"sync"

	"github.com/davinash/yados/internal/events"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreateSqlite() {
	WaitForLeaderElection(suite.cluster)

	for _, s := range suite.cluster.members {
		s.EventHandler().Subscribe(events.CommitEntryEvents)
	}

	defer func() {
		for _, s := range suite.cluster.members {
			s.EventHandler().UnSubscribe(events.CommitEntryEvents)
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
	wg.Wait()

	err = server.ExecuteDDLQuery(&server.QueryArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		StoreName: storeName,
		SQLStr:    "create table employee(empid integer,name varchar(20),title varchar(10))",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
}
