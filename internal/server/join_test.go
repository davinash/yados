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
	defer func(cluster []Server) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

	ports, err := GetFreePorts(1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	newServer, err := StartServerForTests("TestServer-99", "127.0.0.1", int32(ports[0]), nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer StopServer(newServer)

	_, err = cluster[0].AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{Member: newServer.Self()})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	listOfPeers, _ := cluster[0].GetListOfPeers(context.Background(), &pb.ListOfPeersRequest{})
	if len(listOfPeers.GetMember()) != numberOfServers+1 {
		t.Logf("Number of servers did not match. Expected %d Actual %d", numberOfServers+1, len(listOfPeers.GetMember()))
		t.FailNow()
	}

}
