package server

import (
	"context"
	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
	"testing"
)

// TestListAllMembers
// Create 3 node cluster and
// Verify the membership

func verifyListMembers(server *YadosServer, t *testing.T, numberOfServers int) {
	conn, p, err := GetPeerConn(server.self.Address, server.self.Port)
	if err != nil {
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Log(err)
		}
	}(conn)

	peers, err := p.GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(peers.Member) != numberOfServers {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers, len(peers.Member))
		t.FailNow()
	}
}

func TestListAllMembers(t *testing.T) {
	numberOfServers := 3
	cluster, err := CreateClusterForTest(numberOfServers)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []*YadosServer) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

	for _, server := range cluster {
		verifyListMembers(server, t, numberOfServers)
	}
}

func TestListAllMembers2(t *testing.T) {
	numberOfServers := 2
	cluster, err := CreateClusterForTest(numberOfServers)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []*YadosServer) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

	for _, server := range cluster {
		verifyListMembers(server, t, numberOfServers)
	}
}
