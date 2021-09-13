package tests

import (
	"context"
	"fmt"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/raft"
)

func (suite *YadosTestSuite) TestController() {
	_ = WaitForLeaderElection(suite.cluster)

	for _, m := range suite.cluster.members {
		if m.Raft().IsRunning() == false {
			suite.T().Fatalf("Raft on [%s:%s:%d] not started",
				m.Name(), m.Address(), m.Port())
		}
	}

	for _, m := range suite.cluster.members {
		if len(m.Peers()) != len(suite.cluster.members)-1 {
			suite.T().Fatalf("Expected %d, Actual %d", len(suite.cluster.members)-1, len(m.Peers()))
		}
	}

	// Add new member now
	srv, _, err := AddNewServer(3, suite.cluster.members, suite.walDir, "debug", false, suite.controller)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.cluster.members = append(suite.cluster.members, srv)
	time.Sleep(10 * time.Millisecond)

	for _, m := range suite.cluster.members {
		if len(m.Peers()) != len(suite.cluster.members)-1 {
			suite.T().Fatalf("Expected %d, Actual %d", len(suite.cluster.members)-1, len(m.Peers()))
		}
	}

	for _, m := range suite.cluster.members {
		if m.Raft().IsRunning() == false {
			suite.T().Fatalf("Raft on [%s:%s:%d] not started",
				m.Name(), m.Address(), m.Port())
		}
	}
}

func (suite *YadosTestSuite) TestControllerLeader() {
	_ = WaitForLeaderElection(suite.cluster)

	time.Sleep(20 * time.Millisecond)
	leader, err := suite.controller.GetLeader(context.Background(), &pb.GetLeaderRequest{})
	if err != nil {
		return
	}
	for _, server := range suite.cluster.members {
		if server.State() == raft.Leader {
			err := server.Stop()
			if err != nil {
				suite.T().Fatal(err)
			}
			break
		}
	}

	_ = WaitForLeaderElection(suite.cluster)
	for i := 0; i < 5; i++ {
		leader1, err := suite.controller.GetLeader(context.Background(), &pb.GetLeaderRequest{})
		if err != nil {
			return
		}
		fmt.Println(leader.Leader.Name, leader1.Leader.Name)
		if leader.Leader.Name == leader1.Leader.Name && i == 4 {
			suite.T().Fatalf("Leader should have changed")
		}
		time.Sleep(10 * time.Millisecond)
	}

}
