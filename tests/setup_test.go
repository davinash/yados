package tests

import (
	"fmt"
	"github.com/davinash/yados/internal/server"
	"github.com/stretchr/testify/suite"
	"log"
	"sync"
	"testing"
)

type YadosTestSuite struct {
	suite.Suite
	wg sync.WaitGroup
}

func (suite *YadosTestSuite) StartServerForTests(serverName string, port int) {
	defer suite.wg.Done()
	name := fmt.Sprintf("Server-%s", serverName)
	log.Printf("Staring a Server %s\n", name)
	srv, err := server.CreateNewServer(name, "127.0.0.1", port, "Yados Test Cluster")
	if err != nil {
		suite.Failf("Failed to create new server %s, Error = %v", name, err)
	}
	err = srv.Start()
	if err != nil {
		suite.Failf("Failed to create new server %s, Error = %v", name, err)
	}
}

func (suite *YadosTestSuite) SetupTest() {
	log.Println("Setup Test ...")

	suite.wg.Add(1)
	go suite.StartServerForTests("Server-1", 9191)

	suite.wg.Add(1)
	go suite.StartServerForTests("Server-2", 9192)

	suite.wg.Add(1)
	go suite.StartServerForTests("Server-3", 9193)
}

func (suite *YadosTestSuite) Cleanup() {
	log.Println("Cleanup ...")
}

func (suite *YadosTestSuite) TearDownTest() {
	suite.Cleanup()
	suite.wg.Wait()
	log.Println("Tear Down Test ...")
}

func TestYadosTestSuite(t *testing.T) {
	suite.Run(t, new(YadosTestSuite))
}
