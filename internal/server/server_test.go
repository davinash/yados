package server

import (
	"fmt"
	"net"
	"net/http/httptest"
	"testing"
)

func StartServerForTests(name string, address string, port int) (*Server, *httptest.Server, error) {
	server, err := CreateNewServer(name, address, port)
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

	err = server.PostInit(false, "", 0)
	if err != nil {
		return nil, nil, err
	}
	return server, ts, nil
}

func TestCreateNewServer(t *testing.T) {
	server1, httpServer1, err := StartServerForTests("TestServer1", "127.0.0.1", 9090)
	if err != nil {
		t.Logf("Error : %v", err)
	}
	defer httpServer1.Close()

	_, httpServer2, err := StartServerForTests("TestServer2", "127.0.0.1", 9091)
	if err != nil {
		t.Logf("Error : %v", err)
	}
	defer httpServer2.Close()

	_, httpServer3, err := StartServerForTests("TestServer3", "127.0.0.1", 9092)
	if err != nil {
		t.Logf("Error : %v", err)
	}
	defer httpServer3.Close()

	members, err := ListAllMembers(nil, server1)
	if err != nil {
		t.Logf("Error : %v", err)
	}
	t.Logf("%v", members)
}
