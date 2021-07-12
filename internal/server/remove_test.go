package server

import (
	"context"
	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
	"testing"
)

func TestRemoveServer(t *testing.T) {
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
	verifyRemoveServer(t, numberOfServers)
}

func verifyRemoveServer(t *testing.T, numberOfServers int) {
	conn, p, err := GetPeerConn("127.0.0.1", 9191)
	if err != nil {
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Log(err)
		}
	}(conn)

	_, err = p.RemoveServer(context.Background(), &pb.RemoveServerRequest{
		Member: &pb.Member{
			Name:    "TestServer-1",
			Address: "127.0.0.1",
			Port:    9192,
		},
	})
	if err != nil {
		t.FailNow()
	}
	peers, err := p.GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(peers.Member) != numberOfServers-1 {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers, len(peers.Member))
		t.FailNow()
	}

	_, err = p.RemoveServer(context.Background(), &pb.RemoveServerRequest{
		Member: &pb.Member{
			Name:    "TestServer-2",
			Address: "127.0.0.1",
			Port:    9193,
		},
	})
	if err != nil {
		t.FailNow()
	}
	peers, err = p.GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(peers.Member) != numberOfServers-2 {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers, len(peers.Member))
		t.FailNow()
	}
}

func TestRemoveServerNotExists(t *testing.T) {
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

	conn, p, err := GetPeerConn("127.0.0.1", 9191)
	if err != nil {
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Log(err)
		}
	}(conn)

	_, err = p.RemoveServer(context.Background(), &pb.RemoveServerRequest{
		Member: &pb.Member{
			Name:    "TestServer-4",
			Address: "127.0.0.1",
			Port:    9192,
		},
	})
	if err == nil {
		t.FailNow()
	}
}
