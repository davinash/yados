package tests

func (suite *YadosTestSuite) TestStateFileOp() {
	WaitForLeaderElection(suite.cluster)

	err := suite.cluster.members[0].WAL().WriteState(1, "abc")
	if err != nil {
		suite.T().Fatal(err)
	}
	term, votedFor, err := suite.cluster.members[0].WAL().ReadState()
	if err != nil {
		suite.T().Fatal(err)
	}
	if term != 1 {
		suite.T().Fatalf("Failed to read term Expected 1, Actual %d", term)
	}
	if votedFor != "abc" {
		suite.T().Fatalf("Failed to read term Expected abc, Actual %s", votedFor)
	}

	// write new state
	err = suite.cluster.members[0].WAL().WriteState(2, "pqr")
	if err != nil {
		suite.T().Fatal(err)
	}
	term, votedFor, err = suite.cluster.members[0].WAL().ReadState()
	if err != nil {
		suite.T().Fatal(err)
	}
	if term != 2 {
		suite.T().Fatalf("Failed to read term Expected 2, Actual %d", term)
	}
	if votedFor != "pqr" {
		suite.T().Fatalf("Failed to read term Expected pqr, Actual %s", votedFor)
	}
}
