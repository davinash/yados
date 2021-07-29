package server

import (
	"testing"

	"github.com/fortytw2/leaktest"
)

func TestServerStartStop(t *testing.T) {
	defer leaktest.Check(t)()
	cluster, err := CreateNewCluster(3)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()
}
