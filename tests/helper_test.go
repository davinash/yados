package tests

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

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
	logDir  string
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

func AddNewServer(suffix int, members []server.Server, logDir string, logLevel string) (server.Server, error) {
	freePorts, err := GetFreePorts(1)
	if err != nil {
		return nil, err
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
		LogDir:     logDir,
		IsTestMode: true,
	}
	srv, err := server.NewServer(srvArgs)
	if err != nil {
		return nil, err
	}

	err = srv.Serve(peers)
	if err != nil {
		return nil, err
	}
	//suite.cluster.members = append(suite.cluster.members, srv)
	return srv, nil
}

func CreateNewClusterEx(numOfServers int, cluster *TestCluster, logDir string, logLevel string) error {
	srv, err := AddNewServer(0, cluster.members, logDir, logLevel)
	if err != nil {
		panic(err)
	}
	cluster.members = append(cluster.members, srv)

	for i := 1; i < numOfServers; i++ {
		srv, err := AddNewServer(i, cluster.members, logDir, logLevel)
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
	err := CreateNewClusterEx(numOfServers, suite.cluster, suite.logDir, "debug")
	if err != nil {
		return err
	}
	return nil
}

//StopCluster To stop the cluster ( Only for Test Purpose )
func StopCluster(cluster *TestCluster) {
	var wg sync.WaitGroup
	//for i := 0; i < t.numOfServers; i++ {
	//	wg.Add(1)
	//	go func(wg *sync.WaitGroup, index int) {
	//		defer wg.Done()
	//		t.members[index].Stop()
	//	}(&wg, i)
	//}
	for _, srv := range cluster.members {
		wg.Add(1)
		go func(wg *sync.WaitGroup, srv server.Server) {
			defer wg.Done()
			srv.Stop()
		}(&wg, srv)
	}
	wg.Wait()
}

func SetupDataDirectory() string {
	logDir, err := ioutil.TempDir("", "YadosLogStorage")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(logDir), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return logDir
}

func (suite *YadosTestSuite) SetupTest() {
	suite.T().Log("Running SetupTest")

	suite.logDir = SetupDataDirectory()

	err := suite.CreateNewCluster(3)
	if err != nil {
		suite.T().Error(err)
	}

}

func WaitForLeaderElection(cluster *TestCluster) server.Server {
	for _, s := range cluster.members {
		s.EventHandler().Subscribe(server.LeaderChangeEvents)
	}

	defer func() {
		for _, s := range cluster.members {
			s.EventHandler().UnSubscribe(server.LeaderChangeEvents)
		}
	}()

	var set []reflect.SelectCase
	for _, s := range cluster.members {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(s.EventHandler().LeaderChangeEvent()),
		})
	}
	_, valValue, _ := reflect.Select(set)
	peer := valValue.Interface().(server.Server)
	return peer
}

func Cleanup(logDir string) {
	_ = os.RemoveAll(logDir)
}

func (suite *YadosTestSuite) TearDownTest() {
	suite.T().Log("Running TearDownTest")
	StopCluster(suite.cluster)
	Cleanup(suite.logDir)
}

func TestAgentTestSuite(t *testing.T) {
	s := &YadosTestSuite{
		Suite:   suite.Suite{},
		cluster: nil,
	}
	suite.Run(t, s)
}
