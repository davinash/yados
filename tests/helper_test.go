package tests

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
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
func (suite *YadosTestSuite) GetFreePort() (int, *net.TCPListener, error) {
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
func (suite *YadosTestSuite) GetFreePorts(n int) ([]int, error) {
	ports := make([]int, 0)
	for i := 0; i < n; i++ {
		port, listener, err := suite.GetFreePort()
		if err != nil {
			return nil, err
		}
		listener.Close()
		ports = append(ports, port)
	}
	return ports, nil
}

func (suite *YadosTestSuite) AddNewServer(suffix int, isBootStrap bool) error {
	freePorts, err := suite.GetFreePorts(1)
	if err != nil {
		return err
	}
	peers := make([]*pb.Peer, 0)
	for _, p := range suite.cluster.members {
		peers = append(peers, p.Self())
	}
	srvArgs := &server.NewServerArgs{
		Name:     fmt.Sprintf("Server-%d", suffix),
		Address:  "127.0.0.1",
		Port:     int32(freePorts[0]),
		Loglevel: "debug",
		LogDir:   suite.logDir,
	}
	srv, err := server.NewServer(srvArgs)
	if err != nil {
		return err
	}
	err = srv.Serve(peers, isBootStrap)
	if err != nil {
		return err
	}
	suite.cluster.members = append(suite.cluster.members, srv)
	return nil
}

//CreateNewCluster creates a new cluster for the test
func (suite *YadosTestSuite) CreateNewCluster(numOfServers int) error {
	suite.cluster = &TestCluster{
		members:      make([]server.Server, 0),
		numOfServers: numOfServers,
	}
	err := suite.AddNewServer(0, false)
	if err != nil {
		return err
	}
	for i := 1; i < numOfServers; i++ {
		isBootStrap := false
		if i == numOfServers-1 {
			isBootStrap = true
		}
		err := suite.AddNewServer(i, isBootStrap)
		if err != nil {
			return err
		}
	}
	return nil
}

//StopCluster To stop the cluster ( Only for Test Purpose )
func (t *TestCluster) StopCluster() {
	var wg sync.WaitGroup
	for i := 0; i < t.numOfServers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()
			t.members[index].Stop()
		}(&wg, i)
	}
	wg.Wait()
}

func (suite *YadosTestSuite) SetupTest() {
	suite.T().Log("Running SetupTest")

	logDir, err := ioutil.TempDir("", "YadosLogStorage")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(logDir), os.ModePerm)
	if err != nil {
		panic(err)
	}
	suite.logDir = logDir

	err = suite.CreateNewCluster(3)
	if err != nil {
		suite.T().Error(err)
	}
}

func (suite *YadosTestSuite) Cleanup() {
	_ = os.RemoveAll(suite.logDir)
}

func (suite *YadosTestSuite) TearDownTest() {
	suite.T().Log("Running TearDownTest")
	suite.Cleanup()
	suite.cluster.StopCluster()
}

func TestAgentTestSuite(t *testing.T) {
	s := &YadosTestSuite{
		Suite:   suite.Suite{},
		cluster: nil,
	}
	suite.Run(t, s)
}
