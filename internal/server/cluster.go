package server

import (
	"fmt"
	"net"
	"net/http/httptest"
)

type Cluster struct {
}

func startServerForTests(name string, address string, port int) (*Server, *httptest.Server, error) {
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
	return server, ts, nil
}

func addServerInCluster(server *Server, peerAddress string, peerPort int) error {
	err := server.PostInit(true, peerAddress, peerPort)
	if err != nil {
		return err
	}
	return nil
}

func CreateClusterForTest(numOfMembers int) (*Cluster, error) {

	return nil, nil
}

func StopTestCluster(cluster *Cluster) error {
	return nil
}
