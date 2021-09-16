package tests

import (
	"github.com/davinash/yados/internal/raft"
)

func (suite *YadosTestSuite) TestRemoveServer() {
	numOfServers := 7
	for i := 3; i < numOfServers; i++ {
		srv, _, err := AddNewServer(i, suite.cluster.members, suite.walDir, "debug", false,
			suite.controller)
		if err != nil {
			suite.T().Fail()
		}
		suite.cluster.members = append(suite.cluster.members, srv)
	}
	numOfPeers := numOfServers - 1
	for _, peer := range suite.cluster.members {
		if len(peer.Peers()) != numOfPeers {
			suite.T().Fatalf("Number of Peers mismatch, Expected %d, Actual %d for %s", numOfPeers,
				len(peer.Peers()), peer.Name())
		}
	}
	// let us stop a server
	err := suite.cluster.members[0].Stop()
	if err != nil {
		suite.T().Fatal(err)
	}
	numOfPeers--
	for _, peer := range suite.cluster.members {
		if peer.State() == raft.Dead {
			continue
		}
		if len(peer.Peers()) != numOfPeers {
			suite.T().Fatalf("1 :: Number of Peers mismatch, Expected %d, Actual %d for %s", numOfPeers,
				len(peer.Peers()), peer.Name())
		}
	}
}
