package server

import (
	"testing"
)

func TestCreateNewServer(t *testing.T) {
	members, err := CreateClusterForTest(3)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer func(cluster []Cluster) {
		err := StopTestCluster(cluster)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}(members)
}
