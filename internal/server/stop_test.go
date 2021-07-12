package server

import (
	"context"
	"testing"

	pb "github.com/davinash/yados/internal/proto/gen"
)

func TestStopServer(t *testing.T) {
	cluster, err := CreateClusterForTest(3)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for _, server := range cluster {
		_, err := server.StopServer(context.Background(), &pb.StopServerRequest{})
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
}
