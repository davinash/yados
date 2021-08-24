package tests

import (
	"encoding/json"
	"fmt"
	pb "github.com/davinash/yados/internal/proto/gen"
	"io/ioutil"
	"log"
	"net/http"
)

func (suite *YadosTestSuite) TestAPI() {
	WaitForLeaderElection(suite.cluster)

	httpServer, err, port := AddNewServer(len(suite.cluster.members)+1, suite.cluster.members,
		suite.logDir, "debug", true)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.cluster.members = append(suite.cluster.members, httpServer)

	resp, err1 := http.Get(fmt.Sprintf("http://%s:%d/api/v1/status", httpServer.Self().Address, port))
	if err1 != nil {
		suite.T().Fatal(err1)
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer resp.Body.Close()
	sb := string(body)
	var reply pb.ClusterStatusReply
	err = json.Unmarshal([]byte(sb), &reply)
	if err != nil {
		log.Fatal(err)
	}

	suite.T().Log(sb)

}
