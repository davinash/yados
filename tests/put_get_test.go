package tests

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestPutGet() {
	WaitForLeaderElection(suite.cluster)

	_, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: "TestPut",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}

	err = server.ExecuteCmdPut(&server.PutArgs{
		Key:       "Key-1",
		Value:     "Value-1",
		StoreName: "TestPut",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}

	reply, err := server.ExecuteCmdGet(&server.GetArgs{
		Key:       "Key-1",
		StoreName: "TestPut",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	if reply.Value != "Value-1" {
		suite.T().Fatal("Value Mismatch")
	}
}

func (suite *YadosTestSuite) TestPutGetMultiple() {
	WaitForLeaderElection(suite.cluster)

	_, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: "TestPut",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err = server.ExecuteCmdPut(&server.PutArgs{
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: "TestPut",
		}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
		if err != nil {
			suite.T().Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		reply, err := server.ExecuteCmdGet(&server.GetArgs{
			Key:       fmt.Sprintf("Key-%d", i),
			StoreName: "TestPut",
		}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
		if err != nil {
			suite.T().Fatal(err)
		}
		if reply.Value != fmt.Sprintf("Value-%d", i) {
			suite.T().Fatal("Value Mismatch")
		}
	}
}
