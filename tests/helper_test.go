package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/davinash/yados/internal/controller"

	"github.com/davinash/yados/internal/raft"

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
	cluster    *TestCluster
	controller *controller.Controller
	walDir     string
}

//GetFreePort Get the next free port ( Only for test purpose )
func GetFreePort() (int, *net.TCPListener) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	return l.Addr().(*net.TCPAddr).Port, l
}

// GetFreePorts allocates a batch of n TCP ports in one go to avoid collisions. ( Only for test purpose )
func GetFreePorts(n int) []int {
	ports := make([]int, 0)
	for i := 0; i < n; i++ {
		port, listener := GetFreePort()
		err := listener.Close()
		if err != nil {
			panic(err)
		}
		ports = append(ports, port)
	}
	return ports
}

func AddNewServer(suffix int, members []server.Server, walDir string, logLevel string, isWithHTTP bool,
	controller *controller.Controller) (server.Server, int, error) {
	freePorts := GetFreePorts(1)
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
		ports := GetFreePorts(1)
		srvArgs.HTTPPort = ports[0]
	}
	srvArgs.Controllers = make([]string, 0)
	srvArgs.Controllers = append(srvArgs.Controllers, fmt.Sprintf("%s:%d", controller.Address(), controller.Port()))
	srv, err := server.NewServer(srvArgs)
	if err != nil {
		return nil, -1, err
	}
	return srv, srvArgs.HTTPPort, srv.Serve()
}

func CreateNewClusterEx(numOfServers int, cluster *TestCluster, walDir string, logLevel string,
	controller *controller.Controller) error {
	srv, _, err := AddNewServer(0, cluster.members, walDir, logLevel, false, controller)
	if err != nil {
		panic(err)
	}
	cluster.members = append(cluster.members, srv)

	for i := 1; i < numOfServers; i++ {
		srv, _, err := AddNewServer(i, cluster.members, walDir, logLevel, false, controller)
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
	ports := GetFreePorts(1)
	suite.controller = controller.NewController("127.0.0.1", int32(ports[0]), "debug")
	suite.controller.Start()

	err := CreateNewClusterEx(numOfServers, suite.cluster, suite.walDir, "debug", suite.controller)
	if err != nil {
		return err
	}

	return nil
}

//StopCluster To stop the cluster ( Only for Test Purpose )
func StopCluster(cluster *TestCluster, controller *controller.Controller) {
	for _, srv := range cluster.members {
		if srv.State() != raft.Dead {
			err := srv.Stop()
			if err != nil {
				log.Printf("StopCluster -> %v", err)
			}
		}
	}
	controller.Stop()
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

func Cleanup(walDir string) {
	_ = os.RemoveAll(walDir)
}

func (suite *YadosTestSuite) TearDownTest() {
	suite.T().Log("Running TearDownTest")
	StopCluster(suite.cluster, suite.controller)
	Cleanup(suite.walDir)
}

func TestAgentTestSuite(t *testing.T) {
	s := &YadosTestSuite{
		Suite:   suite.Suite{},
		cluster: nil,
	}
	suite.Run(t, s)
}
