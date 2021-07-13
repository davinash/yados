package server

import "testing"

func TestHealthCheck(t *testing.T) {
	numberOfServers := 3
	cluster, err := CreateClusterForTest(numberOfServers)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []*YadosServer) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(cluster)

}
