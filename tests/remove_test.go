package tests

import "github.com/davinash/yados/internal/server"

func (suite *YadosTestSuite) TestRemoveServer() {
	numOfServers := 7
	for i := 4; i < numOfServers; i++ {
		srv, err := AddNewServer(i, suite.cluster.members, suite.logDir, "debug")
		if err != nil {
			suite.T().Fail()
		}
		suite.cluster.members = append(suite.cluster.members, srv)
	}
	numOfPeers := numOfServers - 2
	for _, peer := range suite.cluster.members {
		if len(peer.Peers()) != numOfPeers {
			suite.T().Fatalf("Number of Peers mismatch, Expected %d, Actual %d", numOfPeers,
				len(peer.Peers()))
		}
	}
	// let us stop a server
	err := suite.cluster.members[0].Stop()
	if err != nil {
		suite.T().Fatal(err)
	}
	numOfPeers--
	for _, peer := range suite.cluster.members {
		if peer.State() == server.Dead {
			continue
		}
		if len(peer.Peers()) != numOfPeers {
			suite.T().Fatalf("Number of Peers mismatch, Expected %d, Actual %d", numOfPeers,
				len(peer.Peers()))
		}
	}
}
