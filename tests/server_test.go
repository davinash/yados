package tests

import (
	"time"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestServerStartStop() {
}

func (suite *YadosTestSuite) TestServerLeaderElection() {
	foundLeader := false
	for {
		for _, m := range suite.cluster.members {
			if m.State() == server.Leader {
				foundLeader = true
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func (suite *YadosTestSuite) TestServerLeaderElectionWithWait() {
	foundLeader := false
	var leader server.Server
	for {
		for _, m := range suite.cluster.members {
			if m.State() == server.Leader {
				foundLeader = true
				leader = m
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	counter := 5
	for i := 0; i < counter; i++ {
		suite.T().Log("Waiting and checking for Leader change if any ...")
		for _, m := range suite.cluster.members {
			if m.Name() == leader.Name() {
				if m.State() != server.Leader {
					suite.T().Error("Leader Change not expected")
				}
			}
		}
		//time.Sleep(10 * time.Millisecond)
	}
	suite.T().Log(leader.Name())
}

//func  (suite *YadosTestSuite) TestServerNewLeader(t *testing.T) {
//	defer leaktest.Check(t)()
//
//	foundLeader := false
//	var leader server.Server
//	for {
//		for _, m := range suite.cluster.members {
//			if m.State() == server.Leader {
//				foundLeader = true
//				leader = m
//				break
//			}
//		}
//		if foundLeader {
//			break
//		}
//		time.Sleep(150 * time.Millisecond)
//	}
//	counter := 5
//	for i := 0; i < counter; i++ {
//		t.Log("Waiting and checking for Leader change if any ...")
//		for _, m := range suite.cluster.members {
//			if m.Name() == leader.Name() {
//				if m.State() != server.Leader {
//					t.Error("Leader Change not expected")
//				}
//			}
//		}
//		time.Sleep(150 * time.Millisecond)
//	}
//	// let us stop the leader
//	err := leader.Stop()
//	if err != nil {
//		t.Error(err)
//	}
//	t.Log("Waiting the new leader")
//	foundLeader = false
//	for {
//		for _, m := range suite.cluster.members {
//			if m.State() == server.Leader {
//				foundLeader = true
//				leader = m
//				break
//			}
//		}
//		if foundLeader {
//			break
//		}
//		time.Sleep(150 * time.Millisecond)
//	}
//}
