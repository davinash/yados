package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/davinash/yados/internal/server"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/raft"
)

func (suite *YadosTestSuite) TestController() {
	_, err := server.GetLeader(suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatalf("server.GetLeader failed, error = %v", err)
	}

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
	_, err := server.GetLeader(suite.cluster.members[0].Address(), suite.cluster.members[0].Port())
	if err != nil {
		suite.T().Fatalf("server.GetLeader failed, error = %v", err)
	}

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

TrayAgain:
	leader1, err := suite.controller.GetLeader(context.Background(), &pb.GetLeaderRequest{})
	if err != nil {
		return
	}
	if leader1.Leader == nil {
		time.Sleep(10 * time.Millisecond)
		goto TrayAgain
	}
	fmt.Println(leader.Leader.Name, leader1.Leader.Name)
	if leader.Leader.Name == leader1.Leader.Name {
		suite.T().Fatalf("Leader should have changed")
	}
	time.Sleep(10 * time.Millisecond)
}
