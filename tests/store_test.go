package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/davinash/yados/internal/server"
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

	wg := sync.WaitGroup{}
	for _, p := range suite.cluster.members {
		if len(p.Stores()) != 1 {
			wg.Add(1)
			go func(s server.Server) {
				defer wg.Done()
				storeCreated := false
				for i := 0; i < 5; i++ {
					suite.T().Log("Checking if file exists")
					fileName := filepath.Join(suite.logDir, s.Name(), "log", "testdata", "store.create.TestStoreCreate")
					if _, err := os.Stat(fileName); os.IsNotExist(err) {
						time.Sleep(5 * time.Second)
					} else {
						storeCreated = true
						break
					}
				}
				if !storeCreated {
					suite.T().Errorf("[%s] store count mismatch Exptected 1, found %d", p.Name(), len(p.Stores()))
				}
			}(p)
		}
	}
	wg.Wait()

	for _, p := range suite.cluster.members {
		if _, ok := p.Stores()["TestStoreCreate"]; !ok {
			suite.T().Errorf("store not created on %s", p.Name())
		}
	}
}

func (suite *YadosTestSuite) TestStoreList() {
	suite.WaitForLeaderElection()
	for i := 0; i < 5; i++ {
		err := store.CreateCommandExecute(&store.CreateCommandArgs{
			Address: suite.cluster.members[0].Address(),
			Port:    suite.cluster.members[0].Port(),
			Name:    fmt.Sprintf("TestStoreList-%d", i),
		})
		if err != nil {
			suite.T().Error(err)
		}
	}
	time.Sleep(10 * time.Second)

	storeList, err := store.ExecuteStoreListCommand(&store.ListArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
	})
	if err != nil {
		suite.T().Error(err)
	}
	if len(storeList.Name) != 5 {
		suite.T().Error(err)
	}
}
