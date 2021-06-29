package server

import (
	"fmt"
	"net"
	"net/http/httptest"
	"testing"
)

func StartServerForTests(name string, address string, port int, clusterName string) (*Server, *httptest.Server, error) {
	server, err := CreateNewServer(name, address, port, clusterName)
	if err != nil {
		return nil, nil, err
	}
	server.EnableTestMode()
	testListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, nil, err
	}
	ts := httptest.NewUnstartedServer(SetupRouter(server))
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
	_, httpServer1, err := StartServerForTests("TestServer1", "127.0.0.1", 9090, "TestCluster")
	if err != nil {
		t.Logf("Error : %v", err)
	}
	defer httpServer1.Close()
}
