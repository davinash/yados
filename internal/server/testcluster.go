package server

import (
	"fmt"
	"net"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//GetFreePort Get the next free port ( Only for test purpose )
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

// GetFreePorts allocates a batch of n TCP ports in one go to avoid collisions. ( Only for test purpose )
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

//TestCluster structure to manage the cluster for Tests
type TestCluster struct {
	members      []Server
	numOfServers int
}

//CreateNewCluster creates a new cluster for the test
func CreateNewCluster(numOfServers int) (*TestCluster, error) {
	t := &TestCluster{
		members:      make([]Server, 0),
		numOfServers: numOfServers,
	}
	freePorts, err := GetFreePorts(numOfServers)
	ready := make(chan interface{})

	if err != nil {
		return nil, err
	}

	srv, err := NewServer(fmt.Sprintf("Server-%d", 0), "127.0.0.1",
		int32(freePorts[0]), "debug", ready)
	if err != nil {
		return nil, err
	}
	err = srv.Serve(nil)
	if err != nil {
		return nil, err
	}
	t.members = append(t.members, srv)

	for i := 1; i < numOfServers; i++ {
		peers := make([]*pb.Peer, 0)
		for _, p := range t.members {
			peers = append(peers, p.Self())
		}
		srv, err := NewServer(fmt.Sprintf("Server-%d", i), "127.0.0.1",
			int32(freePorts[i]), "debug", ready)
		if err != nil {
			return nil, err
		}
		err = srv.Serve(peers)
		if err != nil {
			return nil, err
		}
		t.members = append(t.members, srv)
	}
	close(ready)
	return t, nil
}

//StopCluster To stop the cluster ( Only for Test Purpose )
func (t *TestCluster) StopCluster() {
	for i := 0; i < t.numOfServers; i++ {
		t.members[i].RPCServer().Stop()
		t.members[i].Stop()
	}
}
