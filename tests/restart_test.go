package tests

import (
	"fmt"
	"github.com/davinash/yados/internal/controller"

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
	ports := GetFreePorts(1)
	suite.controller = controller.NewController("127.0.0.1", int32(ports[0]), "debug")
	suite.controller.Start()
	err = CreateNewClusterEx(3, suite.cluster, suite.walDir, "debug", suite.controller)
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
