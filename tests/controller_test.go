package tests

func (suite *YadosTestSuite) TestController() {
	_ = WaitForLeaderElection(suite.cluster)

	for _, m := range suite.cluster.members {
		if m.Raft().IsRunning() == false {
			suite.T().Fatalf("Raft on [%s:%s:%d] not started",
				m.Name(), m.Address(), m.Port())
		}
	}
}
