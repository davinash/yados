package tests

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestRestart() {
	WaitForLeaderElection(suite.cluster)
	_, err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Name: "TestRestart",
	}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err = server.ExecuteCmdPut(&server.PutArgs{
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: "TestRestart",
		}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	StopCluster(suite.cluster, suite.controller)

	suite.cluster = &TestCluster{
		members:      make([]server.Server, 0),
		numOfServers: 3,
	}
	err = CreateNewClusterEx(3, suite.cluster, suite.walDir, "debug", nil)
	if err != nil {
		return
	}
	WaitForLeaderElection(suite.cluster)

	for i := 0; i < 10; i++ {
		reply, err := server.ExecuteCmdGet(&server.GetArgs{
			Key:       fmt.Sprintf("Key-%d", i),
			StoreName: "TestRestart",
		}, suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
		if err != nil {
			suite.T().Fatal(err)
		}
		if reply.Value != fmt.Sprintf("Value-%d", i) {
			suite.T().Fatal("Value Mismatch")
		}
	}
}
