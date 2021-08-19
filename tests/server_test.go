package tests

func (suite *YadosTestSuite) TestServerNewLeader() {
	leader := suite.WaitForLeaderElection()

	for i := 4; i < 7; i++ {
		srv, err := AddNewServer(i, suite.cluster.members, suite.logDir)
		if err != nil {
			suite.T().Fail()
		}
		suite.cluster.members = append(suite.cluster.members, srv)
	}
	// let us stop the leader
	err := leader.Stop()
	if err != nil {
		suite.T().Error(err)
	}
	suite.T().Log("Waiting the new leader")

	_ = suite.WaitForLeaderElection()

}
