package server

import (
	"testing"
)

func TestCreateNewServer(t *testing.T) {
	members, err := CreateClusterForTest(3)
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
	}(members)
}

//func TestCreateNewServerDuplicateName(t *testing.T) {
//	cluster, err := CreateClusterForTest(3)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//	defer func(cluster []*server) {
//		err := StopTestCluster(cluster)
//		if err != nil {
//			t.Log(err)
//			t.FailNow()
//		}
//	}(cluster)
//
//	peers := make([]string, 0)
//	for _, p := range cluster {
//		peers = append(peers, fmt.Sprintf("%s:%d", p.self.Address, p.self.Port))
//	}
//	ports, err := GetFreePorts(1)
//	if err != nil {
//		t.Log(err)
//		t.FailNow()
//	}
//	newServer, err := startServerForTests("TestServer-0", "127.0.0.1", int32(ports[0]), peers)
//	if err == nil {
//		t.Log("startServerForTests should fail, duplicate name")
//		t.FailNow()
//	}
//	defer ClearYadosStorage(newServer)
//}

func TestStartGrpcServer(t *testing.T) {
	ports, err := GetFreePorts(1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	srv, err := StartServerForTests("TestServer-0", "127.0.0.1", int32(ports[0]), nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(srv Server) {
		err := StopServer(srv)
		if err != nil {
			t.Logf("Failed to stop the server, error = %v", err)
		}
	}(srv)

	err = srv.Stop()
	if err == nil {
		t.FailNow()
	}
}

func TestPostInitWrongInput(t *testing.T) {
	ports, err := GetFreePorts(1)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	srv, err := StartServerForTests("TestServer-0", "127.0.0.1", int32(ports[0]), nil)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(srv Server) {
		err := StopServer(srv)
		if err != nil {
			t.Logf("Failed to stop the server, error = %v", err)
		}
	}(srv)
	peers := make([]string, 0)
	peers = append(peers, "a:a")

	err = srv.PostInit(peers)
	if err == nil {
		t.Log("PostInit should fail")
		t.FailNow()
	}

	peers1 := make([]string, 0)
	peers1 = append(peers1, "127.0.0.1")

	err = srv.PostInit(peers1)
	if err == nil {
		t.Log("PostInit should fail")
		t.FailNow()
	}
}
