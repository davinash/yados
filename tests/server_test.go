package tests

import (
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"

	"github.com/davinash/yados/internal/server"
)

func (suite *YadosTestSuite) TestServerStartStop() {
}

func (suite *YadosTestSuite) TestServerLeaderElection() {
	foundLeader := false
	for {
		for _, m := range suite.cluster.members {
			if m.State() == pb.Peer_Leader {
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
			if m.State() == pb.Peer_Leader {
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
				if m.State() != pb.Peer_Leader {
					suite.T().Error("Leader Change not expected")
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	suite.T().Log(leader.Name())
}

func (suite *YadosTestSuite) TestServerNewLeader() {
	for i := 4; i < 7; i++ {
		err := suite.AddNewServer(i)
		if err != nil {
			suite.T().Fail()
		}
	}
	foundLeader := false
	var leader server.Server
	for {
		for _, m := range suite.cluster.members {
			if m.State() == pb.Peer_Leader {
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
				if m.State() != pb.Peer_Leader {
					suite.T().Error("Leader Change not expected")
				}
			}
		}
		time.Sleep(150 * time.Millisecond)
	}
	// let us stop the leader
	err := leader.Stop()
	if err != nil {
		suite.T().Error(err)
	}
	suite.T().Log("Waiting the new leader")
	foundLeader = false
	for {
		for _, m := range suite.cluster.members {
			if m.State() == pb.Peer_Leader {
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
}
