package tests

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestPutGet() {
	WaitForLeaderElection(suite.cluster)

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestPut",
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	err = server.ExecuteCmdPut(&server.PutArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		Key:       "Key-1",
		Value:     "Value-1",
		StoreName: "TestPut",
	})
	if err != nil {
		suite.T().Fatal(err)
	}

	reply, err := server.ExecuteCmdGet(&server.GetArgs{
		Address:   suite.cluster.members[0].Address(),
		Port:      suite.cluster.members[0].Port(),
		Key:       "Key-1",
		StoreName: "TestPut",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	if reply.Value != "Value-1" {
		suite.T().Fatal("Value Mismatch")
	}
}

func (suite *YadosTestSuite) TestPutGetMultiple() {
	WaitForLeaderElection(suite.cluster)

	err := server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
		Address: suite.cluster.members[0].Address(),
		Port:    suite.cluster.members[0].Port(),
		Name:    "TestPut",
	})
	if err != nil {
		suite.T().Fatal(err)
	}
	for i := 0; i < 10; i++ {
		err = server.ExecuteCmdPut(&server.PutArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: "TestPut",
		})
		if err != nil {
			suite.T().Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		reply, err := server.ExecuteCmdGet(&server.GetArgs{
			Address:   suite.cluster.members[0].Address(),
			Port:      suite.cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			StoreName: "TestPut",
		})
		if err != nil {
			suite.T().Fatal(err)
		}
		if reply.Value != fmt.Sprintf("Value-%d", i) {
			suite.T().Fatal("Value Mismatch")
		}
	}
}
