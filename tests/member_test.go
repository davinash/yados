package tests

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

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

func (suite *YadosTestSuite) TestClusterStatus() {
	for _, peer := range suite.cluster.members {
		status, err := peer.ClusterStatus(context.Background(), &pb.ClusterStatusRequest{})
		if err != nil {
			suite.T().Error(err)
		}
		if len(status.PeerStatus) != 3 {
			suite.T().Errorf("Number of members mismatch Expected = 3, Actual = %d", len(status.PeerStatus))
		}
	}
}
