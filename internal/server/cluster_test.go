package server

import "fmt"

func startServerForTests(name string, address string, port int32) (*YadosServer, error) {
	server, err := CreateNewServer(name, address, port)
	if err != nil {
		return nil, err
	}
	server.EnableTestMode()
	err = server.StartGrpcServer()
	if err != nil {
		return nil, err
	}
	return server, nil
}

// CreateClusterForTest Creates a cluster for test purpose
func CreateClusterForTest(numOfServers int) ([]*YadosServer, error) {
	portStart := 9191
	cluster := make([]*YadosServer, 0)

	srv, err := startServerForTests("TestServer-0", "127.0.0.1", int32(portStart))
	if err != nil {
		return nil, err
	}
	err = srv.PostInit(nil)
	if err != nil {
		return nil, err
	}

	cluster = append(cluster, srv)
	for i := 1; i < numOfServers; i++ {
		portStart = portStart + 1
		srv, err := startServerForTests(fmt.Sprintf("TestServer-%d", i), "127.0.0.1", int32(portStart))
		if err != nil {
			return nil, err
		}
		peers := make([]string, 0)
		for _, p := range cluster {
			peers = append(peers, fmt.Sprintf("%s:%d", p.self.Address, p.self.Port))
		}
		err = srv.PostInit(peers)
		if err != nil {
			return nil, err
		}
		cluster = append(cluster, srv)
	}
	return cluster, nil
}

//StopTestCluster stops the test cluster
func StopTestCluster(cluster []*YadosServer) error {
	for _, member := range cluster {
		err := member.StopServerFn()
		if err != nil {
			return err
		}
	}
	return nil
}
