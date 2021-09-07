package tests

import (
	"github.com/davinash/yados/internal/raft"
)

func (suite *YadosTestSuite) TestServerNewLeader() {
	_ = WaitForLeaderElection(suite.cluster)

	for i := 4; i < 7; i++ {
		srv, _, err := AddNewServer(i, suite.cluster.members, suite.walDir, "debug	", false)
		if err != nil {
			suite.T().Fail()
		}
		suite.cluster.members = append(suite.cluster.members, srv)
	}

	// let us stop the leader
	for _, server := range suite.cluster.members {
		if server.State() == raft.Leader {
			err := server.Stop()
			if err != nil {
				suite.T().Fatal(err)
			}
			break
		}
	}
	suite.T().Log("Waiting the new leader")

	_ = WaitForLeaderElection(suite.cluster)

}
