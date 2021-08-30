package tests

import (
	"sync"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestStoreCreateSqlite() {
	WaitForLeaderElection(suite.cluster)

	wg := sync.WaitGroup{}
	WaitForEvents(suite.cluster.members, &wg, 2)
	defer func() {
		StopWaitForEvents(suite.cluster.members)
	}()

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
