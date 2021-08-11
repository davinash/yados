package tests

import (
	"github.com/davinash/yados/cmd/cli/commands/store"
)

func (suite *YadosTestSuite) TestStoreCreate() {
	suite.WaitForLeaderElection()
	err := store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestStoreCreate",
	})
	if err != nil {
		suite.T().Error(err)
	}
}
