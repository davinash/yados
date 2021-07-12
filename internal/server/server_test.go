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

	peers := make([]string, 0)
	for _, p := range cluster {
		peers = append(peers, fmt.Sprintf("%s:%d", p.self.Address, p.self.Port))
	}
	_, err = startServerForTests("TestServer-0", "127.0.0.1", 9199, peers)
	if err == nil {
		t.Log("startServerForTests should fail, duplicate name")
		t.FailNow()
	}
}

func TestStartGrpcServer(t *testing.T) {
	srv, err := startServerForTests("TestServer-0", "127.0.0.1", 9191, nil)
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

func TestPostInitWrongInput(t *testing.T) {
	srv, err := startServerForTests("TestServer-0", "127.0.0.1", 9191, nil)
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
	peers := make([]string, 0)
	peers = append(peers, "a:a")

	err = srv.postInit(peers)
	if err == nil {
		t.Log("PostInit should fail")
		t.FailNow()
	}

	peers1 := make([]string, 0)
	peers1 = append(peers1, "127.0.0.1")

	err = srv.postInit(peers1)
	if err == nil {
		t.Log("PostInit should fail")
		t.FailNow()
	}
}
