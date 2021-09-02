package tests

import (
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/wal"
	"github.com/google/uuid"
)

func (suite *YadosTestSuite) TestWALIterator() {
	WaitForLeaderElection(suite.cluster)
	numOfEntries := 100
	ids := make([]string, 0)
	for i := 0; i < numOfEntries; i++ {
		id := uuid.New().String()
		err := suite.cluster.members[0].WAL().Append(&pb.WalEntry{
			Term:    0,
			Index:   0,
			Command: nil,
			Id:      id,
			CmdType: 0,
		})
		if err != nil {
			suite.T().Fatal(err)
		}
		ids = append(ids, id)
	}
	iter, err := suite.cluster.members[0].WAL().Iterator()
	if err != nil {
		suite.T().Fatal(err)
	}

	defer func(iter wal.Iterator) {
		err := iter.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}(iter)

	entry, err := iter.Next()
	if err != nil {
		return
	}

	index := 0
	for entry != nil {
		if entry.Id != ids[index] {
			suite.T().Fatalf("Entries mismatch Expected Id %s, Actual Id = %s", ids[index], entry.Id)
		}

		entry, err = iter.Next()
		if err != nil {
			return
		}
		index++
	}
	if index != numOfEntries {
		suite.T().Fatalf("Number of Entries mismatch Expected %d, Actual %d", index, numOfEntries)
	}

}

func (suite *YadosTestSuite) TestWALIteratorEmpty() {
	WaitForLeaderElection(suite.cluster)

	iter, err := suite.cluster.members[0].WAL().Iterator()
	if err != nil {
		suite.T().Fatal(err)
	}
	defer func(iter wal.Iterator) {
		err := iter.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}(iter)

	entry, err1 := iter.Next()
	if err1 != nil {
		return
	}

	count := 0
	for entry != nil {
		count++
		entry, err = iter.Next()
		if err != nil {
			return
		}
	}
	if count != 0 {
		suite.T().Fatalf("Number of Entries mismatch Expected 0, Actual %d", count)
	}
}
