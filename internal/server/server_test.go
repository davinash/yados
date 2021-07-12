package server

import (
	"fmt"
	"testing"
)

func TestCreateNewServer(t *testing.T) {
	members, err := CreateClusterForTest(3)
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
	}(members)
}

func TestCreateNewServerDuplicateName(t *testing.T) {
	cluster, err := CreateClusterForTest(3)
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

	srvDuplicate, err := startServerForTests("TestServer-0", "127.0.0.1", 9199)
	if err != nil {
		t.FailNow()
	}
	peers := make([]string, 0)
	for _, p := range cluster {
		peers = append(peers, fmt.Sprintf("%s:%d", p.self.Address, p.self.Port))
	}
	err = srvDuplicate.PostInit(peers)
	if err == nil {
		t.Log("startServerForTests should fail, duplicate name")
		t.FailNow()
	}
}

func TestStartGrpcServer(t *testing.T) {
	srv, err := startServerForTests("TestServer-0", "127.0.0.1", 9191)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	err = srv.PostInit(nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(srv *YadosServer) {
		err := srv.StopServerFn()
		if err != nil {
			t.Logf("Failed to stop the server, error = %v", err)
		}
	}(srv)

	err = srv.StartGrpcServer()
	if err == nil {
		t.FailNow()
	}
}
