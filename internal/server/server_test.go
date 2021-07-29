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

func TestServerLeaderElectionWithWait(t *testing.T) {
	defer leaktest.Check(t)()
	cluster, err := CreateNewCluster(3)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()

	foundLeader := false
	var leader Server
	for {
		for _, m := range cluster.members {
			if m.State() == Leader {
				foundLeader = true
				leader = m
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	counter := 5
	for i := 0; i < counter; i++ {
		t.Log("Waiting and checking for Leader change if any ...")
		for _, m := range cluster.members {
			if m.Name() == leader.Name() {
				if m.State() != Leader {
					t.Error("Leader Change not expected")
				}
			}
		}
		time.Sleep(150 * time.Millisecond)
	}
	t.Log(leader.Name())
}

func TestServerNewLeader(t *testing.T) {
	defer leaktest.Check(t)()
	cluster, err := CreateNewCluster(7)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()

	foundLeader := false
	var leader Server
	for {
		for _, m := range cluster.members {
			if m.State() == Leader {
				foundLeader = true
				leader = m
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	counter := 5
	for i := 0; i < counter; i++ {
		t.Log("Waiting and checking for Leader change if any ...")
		for _, m := range cluster.members {
			if m.Name() == leader.Name() {
				if m.State() != Leader {
					t.Error("Leader Change not expected")
				}
			}
		}
		time.Sleep(150 * time.Millisecond)
	}
	// let us stop the leader
	err = leader.Stop()
	if err != nil {
		t.Error(err)
	}
	t.Log("Waiting the new leader")
	foundLeader = false
	for {
		for _, m := range cluster.members {
			if m.State() == Leader {
				foundLeader = true
				leader = m
				break
			}
		}
		if foundLeader {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
}
