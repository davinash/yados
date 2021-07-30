package tests

func (suite *YadosTestSuite) TestAddNewMember() {
	numOfMembers := 3
	for _, peer := range suite.cluster.members {
		actMembers := len(peer.Peers())
		if actMembers != numOfMembers-1 {
			suite.T().Errorf("member count mismatch for server %s, Expected %d, Actual %d", peer.Name(),
				numOfMembers-1, actMembers)
		}
	}
}
