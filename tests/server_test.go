package tests

func (suite *YadosTestSuite) TestServerNewLeader() {
	leader := suite.WaitForLeaderElection()

	for i := 4; i < 7; i++ {
		err := suite.AddNewServer(i)
		if err != nil {
			suite.T().Fail()
		}
	}
	// let us stop the leader
	err := leader.Stop()
	if err != nil {
		suite.T().Error(err)
	}
	suite.T().Log("Waiting the new leader")

	_ = suite.WaitForLeaderElection()

}
