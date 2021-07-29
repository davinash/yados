package server

import (
	"testing"
)

func TestServerStartStop(t *testing.T) {
	//defer leaktest.Check(t)()
	cluster, err := CreateNewCluster(3)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()
}
