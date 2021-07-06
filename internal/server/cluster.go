package server

import (
	"fmt"
	"net"
	"net/http/httptest"
)

//Cluster data structure to represent cluster. Only for test purpose
type Cluster struct {
	Member  *Server
	HTTPSrv *httptest.Server
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

// CreateClusterForTest Creates a cluster for test purpose
func CreateClusterForTest(numOfMembers int) ([]Cluster, error) {
	portStart := 9191
	cluster := make([]Cluster, 0)

	srv, httpSrv, err := startServerForTests("TestServer-0", "127.0.0.1", portStart)
	if err != nil {
		return nil, err
	}
	cluster = append(cluster, Cluster{
		Member:  srv,
		HTTPSrv: httpSrv,
	})

	for i := 1; i < numOfMembers; i++ {
		portStart = portStart + 1
		srv, httpSrv, err := startServerForTests(fmt.Sprintf("TestServer-%d", i), "127.0.0.1", portStart)
		if err != nil {
			return nil, err
		}
		err = srv.PostInit(true, cluster[i-1].Member.self.Address, cluster[i-1].Member.self.Port)
		if err != nil {
			return nil, err
		}
		cluster = append(cluster, Cluster{
			Member:  srv,
			HTTPSrv: httpSrv,
		})
	}
	return cluster, nil
}

//StopTestCluster stops the test cluster
func StopTestCluster(cluster []Cluster) error {
	for _, member := range cluster {
		err := member.Member.StopServerFn()
		if err != nil {
			return err
		}
		member.HTTPSrv.Close()
	}
	return nil
}
