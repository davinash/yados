package server

import (
	"context"
	"testing"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
)

func TestRaft_Start_Stop(t *testing.T) {
	numOfServers := 1

	cluster, err := CreateClusterForTest(numOfServers)
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
}

func TestRaft_LeaderElection(t *testing.T) {
	numOfServers := 3
	storeName := "Store-1"

	cluster, err := CreateClusterForTest(numOfServers)
	if err != nil {
		t.Error(err)
	}
	defer func(cluster []Server) {
		t.Log("Stopping the cluster now")
		err := StopTestCluster(cluster)
		if err != nil {
			t.Error(err)
		}
	}(cluster)

	_, err = cluster[0].CreateStore(context.Background(), &pb.StoreCreateRequest{
		Name:        storeName,
		Replication: int32(numOfServers) - 1,
	})
	if err != nil {
		t.Error(err)
	}
	for {
		if Leader == cluster[0].Stores()[storeName].RaftInstance().State() ||
			Leader == cluster[1].Stores()[storeName].RaftInstance().State() ||
			Leader == cluster[2].Stores()[storeName].RaftInstance().State() {
			break
		}
		time.Sleep(DefaultElectionTimeout)
	}
}
