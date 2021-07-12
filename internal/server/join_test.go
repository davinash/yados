package server

import (
	"context"
	"testing"

	pb "github.com/davinash/yados/internal/proto/gen"
)

func TestAddNewMemberInCluster(t *testing.T) {
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

	newServer, err := startServerForTests("TestServer-99", "127.0.0.1", 9199, nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer newServer.StopServerFn()

	_, err = cluster[0].AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{Member: newServer.self})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	listOfPeers := cluster[0].GetListOfPeersEx()
	if len(listOfPeers) != numberOfServers+1 {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers+1, len(listOfPeers))
		t.FailNow()
	}

}
