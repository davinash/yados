package server

import (
	"net"
	"net/http/httptest"
	"testing"
)

func StartServerForTests(name string, address string, port int, clusterName string, t *testing.T) (*Server, *httptest.Server, error) {
	server, err := CreateNewServer(name, address, port, clusterName)
	if err != nil {
		t.Logf("Error = %v", err)
		t.FailNow()
	}

	testListener, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		return nil, nil, err
	}
	ts := httptest.NewUnstartedServer(nil)
	ts.Listener.Close()
	ts.Listener = testListener
	ts.Start()

	err = server.PostInit()
	if err != nil {
		return nil, nil, err
	}
	return server, ts, nil
}

func TestCreateNewServer(t *testing.T) {
	_, httpServer1, err := StartServerForTests("TestServer1", "127.0.0.1", 9090, "TestCluster", t)
	if err != nil {
		t.Logf("Error : %v", err)
	}
	defer httpServer1.Close()

}
