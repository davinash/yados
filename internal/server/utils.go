package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//GetPeerConn returns the connection and grpc client for the remote peer
func GetPeerConn(address string, port int32) (*grpc.ClientConn, pb.YadosServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	peer := pb.NewYadosServiceClient(conn)
	return conn, peer, nil
}

//StartServerForTests starts the servers for the tests ( Only for test purpose )
func StartServerForTests(name string, address string, port int32, peers []string) (Server, error) {
	freePorts, err := GetFreePorts(1)
	if err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir("", "YadosStorage")
	if err != nil {
		return nil, err
	}

	server, err := CreateNewServer(&SrvArgs{
		Name:    name,
		Address: address,
		Port:    int32(freePorts[0]),
		HomeDir: dir,
	})
	if err != nil {
		return nil, err
	}
	server.EnableTestMode()
	server.SetLogLevel("debug")

	err = server.Start(peers)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// CreateClusterForTest Creates a cluster for test purpose ( Only for test purpose )
func CreateClusterForTest(numOfServers int) ([]Server, error) {
	freePorts, err := GetFreePorts(numOfServers)
	if err != nil {
		return nil, err
	}
	cluster := make([]Server, 0)

	srv, err := StartServerForTests("TestServer-0", "127.0.0.1", int32(freePorts[0]), nil)
	if err != nil {
		return nil, err
	}

	cluster = append(cluster, srv)
	for i := 1; i < numOfServers; i++ {
		peers := make([]string, 0)
		for _, p := range cluster {
			peers = append(peers, fmt.Sprintf("%s:%d", p.Address(), p.Port()))
		}
		srv, err := StartServerForTests(fmt.Sprintf("TestServer-%d", i), "127.0.0.1", int32(freePorts[i]), peers)
		if err != nil {
			return nil, err
		}

		cluster = append(cluster, srv)
	}
	return cluster, nil
}

//ClearYadosStorage deletes the temporary storage created ( Only for test purpose )
func ClearYadosStorage(server Server) error {
	err := os.RemoveAll(server.DataDirectory())
	if err != nil {
		return err
	}
	return nil
}

//StopServer stops the servers used in tests ( Only for test purpose )
func StopServer(server Server) error {
	err := server.Stop()
	if err != nil {
		return err
	}
	err = ClearYadosStorage(server)
	if err != nil {
		return err
	}
	return nil
}

//StopTestCluster stops the test cluster ( Only for test purpose )
func StopTestCluster(cluster []Server) error {
	for _, member := range cluster {
		err := StopServer(member)
		if err != nil {
			return err
		}
	}
	return nil
}

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
