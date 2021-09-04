package tests

import (
	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestDeleteStore() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestDeleteStore"
	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	listStore, err := server.ExecuteCmdListStore(&server.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	if len(listStore.Stores) != 1 {
		suite.T().Fatalf("Store count mismatch Expected 1, Actual %d", len(listStore.Stores))
	}

	err = store.ExecuteCmdDeleteStore(&store.DeleteArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		StoreName: storeName,
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	if len(listStore.Stores) == 0 {
		suite.T().Fatalf("Store count mismatch Expected 0, Actual %d", len(listStore.Stores))
	}
}
