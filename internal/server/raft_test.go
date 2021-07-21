package server

import (
	"context"
	"testing"

	pb "github.com/davinash/yados/internal/proto/gen"
)

func TestRaft_Start(t *testing.T) {
	cluster, err := CreateClusterForTest(3)
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

	_, err = cluster[0].CreateStore(context.Background(), &pb.StoreCreateRequest{
		Name:        "Store-1",
		Replication: 3,
	})
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, p := range cluster {
		if len(p.Stores()) != 1 {
			t.FailNow()
		}
	}
}
