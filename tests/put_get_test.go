package tests

import (
	"fmt"

	"github.com/davinash/yados/cmd/cli/commands/store"
)

func (suite *YadosTestSuite) TestPutGet() {
	suite.WaitForLeaderElection()

	err := store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestPut",
	})
	if err != nil {
		suite.T().Error(err)
	}

	err = store.ExecutePutCommand(&store.PutArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		Key:       "Key-1",
		Value:     "Value-1",
		StoreName: "TestPut",
	})
	if err != nil {
		suite.T().Error(err)
	}

	reply, err := store.ExecuteGetCommand(&store.GetArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		Key:       "Key-1",
		StoreName: "TestPut",
	})
	if err != nil {
		suite.T().Error(err)
	}
	if reply.Value != "Value-1" {
		suite.T().Error("Value Mismatch")
	}
}

func (suite *YadosTestSuite) TestPutGetMultiple() {
	suite.WaitForLeaderElection()

	err := store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestPut",
	})
	if err != nil {
		suite.T().Error(err)
	}
	for i := 0; i < 10; i++ {
		err = store.ExecutePutCommand(&store.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: "TestPut",
		})
		if err != nil {
			suite.T().Error(err)
		}
	}

	for i := 0; i < 10; i++ {
		reply, err := store.ExecuteGetCommand(&store.GetArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			StoreName: "TestPut",
		})
		if err != nil {
			suite.T().Error(err)
		}
		if reply.Value != fmt.Sprintf("Value-%d", i) {
			suite.T().Error("Value Mismatch")
		}
	}
}
