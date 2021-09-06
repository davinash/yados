package tests

import (
	"sync"

	"github.com/davinash/yados/internal/server"
)

func CreateSQLStoreAndWait(cluster *TestCluster, storeName string) error {
	wg := sync.WaitGroup{}
	WaitForEvents(cluster.members, &wg, 1)
	defer func() {
		StopWaitForEvents(cluster.members)
	}()
	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: storeName,
		Type: "sqlite",
	}, cluster.members[0].Address(), cluster.members[0].Port())
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func CreateTableAndWait(cluster *TestCluster, storeName string) error {
	wg := sync.WaitGroup{}
	WaitForEvents(cluster.members, &wg, 1)
	defer func() {
		StopWaitForEvents(cluster.members)
	}()

	_, err := server.ExecuteCmdQuery(&server.QueryArgs{
		StoreName: storeName,
		SQLStr:    "create table employee(empid integer,name varchar(20),title varchar(10))",
	}, cluster.members[0].Address(), cluster.members[0].Port())
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func InsertRowsAndWait(cluster *TestCluster, storeName string) error {
	wg := sync.WaitGroup{}
	WaitForEvents(cluster.members, &wg, 5)
	defer func() {
		StopWaitForEvents(cluster.members)
	}()

	queries := []string{
		"insert into employee values(101,'John Smith','CEO')",
		"insert into employee values(102,'Raj Reddy','Sysadmin')",
		"insert into employee values(103,'Jason Bourne','Developer')",
		"insert into employee values(104,'Jane Smith','Sale Manager')",
		"insert into employee values(105,'Rita Patel','DBA')",
	}
	for _, q := range queries {
		_, err := server.ExecuteCmdQuery(&server.QueryArgs{
			StoreName: storeName,
			SQLStr:    q,
		}, cluster.members[0].Address(), cluster.members[0].Port())
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func (suite *YadosTestSuite) TestStoreCreateSqlite() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestStoreCreateSqlite"

	if err := CreateSQLStoreAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := CreateTableAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := InsertRowsAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *YadosTestSuite) TestStoreSelectSqlite() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestStoreSelectSqlite"

	if err := CreateSQLStoreAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := CreateTableAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := InsertRowsAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	reply, err := server.ExecuteCmdSQLQuery(&server.QueryArgs{
		StoreName: storeName,
		SQLStr:    "SELECT * from employee",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}

	if len(reply.Rows) != 5 {
		suite.T().Fatalf("expected Rows = 5, Actual = %d", len(reply.Rows))
	}
}
