package server

import (
	"fmt"
	"net"
)

func startServerForTests(name string, address string, port int32, peers []string) (*YadosServer, error) {
	freePorts, err := GetFreePorts(1)
	if err != nil {
		return nil, err
	}
	server, err := CreateNewServer(name, address, int32(freePorts[0]))
	if err != nil {
		return nil, err
	}
	server.EnableTestMode()
	err = server.Start(peers)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// CreateClusterForTest Creates a cluster for test purpose
func CreateClusterForTest(numOfServers int) ([]*YadosServer, error) {
	freePorts, err := GetFreePorts(numOfServers)
	if err != nil {
		return nil, err
	}
	cluster := make([]*YadosServer, 0)

	srv, err := startServerForTests("TestServer-0", "127.0.0.1", int32(freePorts[0]), nil)
	if err != nil {
		return nil, err
	}

	cluster = append(cluster, srv)
	for i := 1; i < numOfServers; i++ {
		peers := make([]string, 0)
		for _, p := range cluster {
			peers = append(peers, fmt.Sprintf("%s:%d", p.self.Address, p.self.Port))
		}
		srv, err := startServerForTests(fmt.Sprintf("TestServer-%d", i), "127.0.0.1", int32(freePorts[i]), peers)
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

func GetFreePort() (int, *net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, nil, err
	}
	return l.Addr().(*net.TCPAddr).Port, l, nil
}

// GetFreePorts allocates a batch of n TCP ports in one go to avoid collisions.
func GetFreePorts(n int) ([]int, error) {
	ports := make([]int, 0)
	for i := 0; i < n; i++ {
		port, listener, err := GetFreePort()
		if err != nil {
			return nil, err
		}
		listener.Close()
		ports = append(ports, port)
	}
	return ports, nil
}
