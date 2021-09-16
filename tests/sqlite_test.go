package tests

import (
	"fmt"
	"log"
	"sync"

	"github.com/davinash/yados/internal/server"
)

func CreateTableAndWaitFail(cluster *TestCluster, storeName string) error {
	wg := sync.WaitGroup{}
	WaitForEvents(cluster.members, &wg, 1)
	defer func() {
		StopWaitForEvents(cluster.members)
	}()

	_, err := server.ExecuteCmdQuery(&server.QueryArgs{
		StoreName: storeName,
		SQLStr:    "create table employee(empid integer,name varchar(20)title varchar(10))",
	}, cluster.members[0].Address(), cluster.members[0].Port())
	if err == nil {
		return fmt.Errorf("should have thrown error")
	}
	log.Println(err.Error())
	wg.Wait()
	return nil
}

func (suite *YadosTestSuite) TestStoreSqlite_InvalidSQL() {
	storeName := "TestStoreSqlite_InvalidSQL"
	if err := CreateSQLStoreAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := CreateTableAndWaitFail(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
}
