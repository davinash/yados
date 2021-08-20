package tests

import "github.com/davinash/yados/cmd/cli/commands/store"

func (suite *YadosTestSuite) TestDeleteStore() {
	WaitForLeaderElection(suite.cluster)
	storeName := "TestDeleteStore"
	err := store.ExecuteCmdCreateStore(&store.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		suite.T().Error(err)
	}
	listStore, err := store.ExecuteCmdListStore(&store.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		return
	}

	if len(listStore.Name) != 1 {
		suite.T().Errorf("Store count mismatch Expected 1, Actual %d", len(listStore.Name))
	}

	err = store.ExecuteCmdDeleteStore(&store.DeleteArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		StoreName: storeName,
	})
	if err != nil {
		return
	}

	if len(listStore.Name) == 0 {
		suite.T().Errorf("Store count mismatch Expected 0, Actual %d", len(listStore.Name))
	}
}
