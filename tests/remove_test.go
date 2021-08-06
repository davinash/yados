package tests

import "github.com/davinash/yados/internal/server"

func (suite *YadosTestSuite) TestRemoveServer() {
	numOfServers := 7
	for i := 4; i < numOfServers; i++ {
		err := suite.AddNewServer(i)
		if err != nil {
			suite.T().Fail()
		}
	}
	numOfPeers := numOfServers - 2
	for _, peer := range suite.cluster.members {
		if len(peer.Peers()) != numOfPeers {
			suite.T().Errorf("Number of Peers mismatch, Expected %d, Actual %d", numOfPeers,
				len(peer.Peers()))
		}
	}
	// let us stop a server
	err := suite.cluster.members[0].Stop()
	if err != nil {
		suite.T().Error(err)
	}
	numOfPeers--
	for _, peer := range suite.cluster.members {
		if peer.State() == server.Dead {
			continue
		}
		if len(peer.Peers()) != numOfPeers {
			suite.T().Errorf("Number of Peers mismatch, Expected %d, Actual %d", numOfPeers,
				len(peer.Peers()))
		}
	}
}
