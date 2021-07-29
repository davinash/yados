package server

import (
	"testing"

	"github.com/fortytw2/leaktest"
)

func TestAddNewMember(t *testing.T) {
	defer leaktest.Check(t)()
	numOfMembers := 3
	cluster, err := CreateNewCluster(numOfMembers)
	if err != nil {
		t.Error(err)
	}
	defer cluster.StopCluster()

	for _, peer := range cluster.members {
		actMembers := len(peer.Peers())
		if actMembers != numOfMembers-1 {
			t.Errorf("member count mismatch for server %s, Expected %d, Actual %d", peer.Name(),
				numOfMembers-1, actMembers)
		}
	}
}
