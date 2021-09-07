package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/davinash/yados/internal/raft"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/stretchr/testify/suite"
)

//TestCluster structure to manage the cluster for Tests
type TestCluster struct {
	members      []server.Server
	numOfServers int
}

type YadosTestSuite struct {
	suite.Suite
	cluster *TestCluster
	walDir  string
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

func AddNewServer(suffix int, members []server.Server, walDir string, logLevel string, isWithHTTP bool) (server.Server, int, error) {
	freePorts, err := GetFreePorts(1)
	if err != nil {
		return nil, -1, err
	}
	peers := make([]*pb.Peer, 0)
	for _, p := range members {
		peers = append(peers, p.Self())
	}
	srvArgs := &server.NewServerArgs{
		Name:       fmt.Sprintf("Server-%d", suffix),
		Address:    "127.0.0.1",
		Port:       int32(freePorts[0]),
		Loglevel:   logLevel,
		WalDir:     walDir,
		IsTestMode: true,
		HTTPPort:   -1,
	}

	if isWithHTTP {
		ports, err := GetFreePorts(1)
		if err != nil {
			return nil, -1, err
		}
		srvArgs.HTTPPort = ports[0]
	}
	srv, err := server.NewServer(srvArgs)
	if err != nil {
		return nil, -1, err
	}

	err = srv.Serve(peers)
	if err != nil {
		return nil, -1, err
	}
	return srv, srvArgs.HTTPPort, nil
}

func CreateNewClusterEx(numOfServers int, cluster *TestCluster, walDir string, logLevel string) error {
	srv, _, err := AddNewServer(0, cluster.members, walDir, logLevel, false)
	if err != nil {
		panic(err)
	}
	cluster.members = append(cluster.members, srv)

	for i := 1; i < numOfServers; i++ {
		srv, _, err := AddNewServer(i, cluster.members, walDir, logLevel, false)
		if err != nil {
			panic(err)
		}
		cluster.members = append(cluster.members, srv)
	}
	return nil
}

//CreateNewCluster creates a new cluster for the test
func (suite *YadosTestSuite) CreateNewCluster(numOfServers int) error {
	suite.cluster = &TestCluster{
		members:      make([]server.Server, 0),
		numOfServers: numOfServers,
	}
	err := CreateNewClusterEx(numOfServers, suite.cluster, suite.walDir, "debug")
	if err != nil {
		return err
	}
	return nil
}

//StopCluster To stop the cluster ( Only for Test Purpose )
func StopCluster(cluster *TestCluster) {
	var wg sync.WaitGroup
	for _, srv := range cluster.members {
		wg.Add(1)
		go func(wg *sync.WaitGroup, srv server.Server) {
			defer wg.Done()
			if srv.State() != raft.Dead {
				err := srv.Stop()
				if err != nil {
					log.Printf("StopCluster -> %v", err)
				}
			}
		}(&wg, srv)
	}
	wg.Wait()
}

func SetupDataDirectory() string {
	walDir, err := ioutil.TempDir("", "YadosWALStorage")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(walDir), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return walDir
}

func (suite *YadosTestSuite) SetupTest() {
	suite.T().Log("Running SetupTest")

	suite.walDir = SetupDataDirectory()

	err := suite.CreateNewCluster(3)
	if err != nil {
		suite.T().Fatal(err)
	}

}

func WaitForLeaderElection(cluster *TestCluster) *pb.Peer {
	for _, s := range cluster.members {
		s.EventHandler().LeaderChangeChan = make(chan interface{})
	}
	defer func() {
		for _, s := range cluster.members {
			close(s.EventHandler().LeaderChangeChan)
			s.EventHandler().LeaderChangeChan = nil
		}
	}()

	var set []reflect.SelectCase
	for _, s := range cluster.members {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(s.EventHandler().LeaderChangeChan),
		})
	}
	_, valValue, _ := reflect.Select(set)
	peer := valValue.Interface().(*pb.Peer)
	return peer
}

func Cleanup(walDir string) {
	_ = os.RemoveAll(walDir)
}

func (suite *YadosTestSuite) TearDownTest() {
	suite.T().Log("Running TearDownTest")
	StopCluster(suite.cluster)
	Cleanup(suite.walDir)
}

func TestAgentTestSuite(t *testing.T) {
	s := &YadosTestSuite{
		Suite:   suite.Suite{},
		cluster: nil,
	}
	suite.Run(t, s)
}
