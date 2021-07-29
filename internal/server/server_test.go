package server

import (
	"testing"
	"time"

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

func TestServerLeaderElection(t *testing.T) {
	defer leaktest.Check(t)()
	cluster, err := CreateNewCluster(3)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()

	foundLeader := false
	for {
		for _, m := range cluster.members {
			if m.State() == Leader {
				foundLeader = true
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
}
