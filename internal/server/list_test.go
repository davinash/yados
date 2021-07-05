package server

import (
	"testing"
)

// TestListAllMembers
// Create 3 node cluster and
// Verify the membership

func TestListAllMembers(t *testing.T) {
	//members, err := CreateClusterForTest(3)
	//if err != nil {
	//	t.Log(err)
	//	t.FailNow()
	//}
	//defer func(cluster []Cluster) {
	//	err := StopTestCluster(cluster)
	//	if err != nil {
	//		t.Log(err)
	//		t.FailNow()
	//	}
	//}(members)
	//
	//for _, member := range members {
	//	m, err := ListAllMembers(nil, member.Member)
	//	if err != nil {
	//		t.Log(err)
	//		t.FailNow()
	//	}
	//	v, ok := m.Resp.([]MemberServer)
	//	if !ok {
	//		t.FailNow()
	//	}
	//	fmt.Println(v)
	//	if len(v) != len(members) {
	//		t.Logf("Cluster Members did not match, Expected %d , Actual %d", len(members), len(v))
	//		t.FailNow()
	//	}
	//}
}
