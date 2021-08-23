package tests

import (
	"fmt"

	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestRestart() {
	WaitForLeaderElection(suite.cluster)
	err := store.ExecuteCmdCreateStore(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestRestart",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err = store.ExecuteCmdPut(&store.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: "TestRestart",
		})
		if err != nil {
			suite.T().Fatal(err)
		}
	}
	StopCluster(suite.cluster)

	suite.cluster = &TestCluster{
		members:      make([]server.Server, 0),
		numOfServers: 3,
	}
	err = CreateNewClusterEx(3, suite.cluster, suite.logDir, "debug")
	if err != nil {
		return
	}
	WaitForLeaderElection(suite.cluster)

	for i := 0; i < 10; i++ {
		reply, err := store.ExecuteCmdGet(&store.GetArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			StoreName: "TestRestart",
		})
		if err != nil {
			suite.T().Fatal(err)
		}
		if reply.Value != fmt.Sprintf("Value-%d", i) {
			suite.T().Fatal("Value Mismatch")
		}
	}
}
