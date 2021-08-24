package tests

func (suite *YadosTestSuite) TestServerNewLeader() {
	leader := WaitForLeaderElection(suite.cluster)

	for i := 4; i < 7; i++ {
		srv, err, _ := AddNewServer(i, suite.cluster.members, suite.logDir, "debug", false)
		if err != nil {
			suite.T().Fail()
		}
		suite.cluster.members = append(suite.cluster.members, srv)
	}
	// let us stop the leader
	err := leader.Stop()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.T().Log("Waiting the new leader")

	_ = WaitForLeaderElection(suite.cluster)

}
