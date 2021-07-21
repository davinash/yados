package server

import (
	"context"
	"testing"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

func TestRemoveServer(t *testing.T) {
	numberOfServers := 3
	cluster, err := CreateClusterForTest(numberOfServers)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []Server) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)
	verifyRemoveServer(t, numberOfServers, cluster)
}

func verifyRemoveServer(t *testing.T, numberOfServers int, cluster []Server) {
	conn, p, err := GetPeerConn("127.0.0.1", cluster[0].Self().Port)
	if err != nil {
		t.Log(err)
		t.FailNow()
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
			Port:    cluster[1].Self().Port,
		},
	})
	if err != nil {
		t.Log(err)
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
			Port:    cluster[2].Self().Port,
		},
	})
	if err != nil {
		t.Log(err)
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
	defer func(cluster []Server) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

	conn, p, err := GetPeerConn("127.0.0.1", cluster[0].Self().Port)
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
			Port:    cluster[1].Self().Port,
		},
	})
	if err == nil {
		t.FailNow()
	}
}

func TestRemoveServerFromCluster(t *testing.T) {
	numberOfServers := 3
	cluster, err := CreateClusterForTest(numberOfServers)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []Server) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

	err = RemoveServerFromCluster(cluster[0], "127.0.0.1", cluster[1].Self().Port)
	if err != nil {
		t.Logf("Error %v", err)
		t.FailNow()
	}
	listOfPeersEx, _ := cluster[1].GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
	if len(listOfPeersEx.GetMember()) != numberOfServers-1 {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers-1,
			len(listOfPeersEx.GetMember()))
		t.FailNow()
	}
}
