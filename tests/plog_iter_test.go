package tests

import (
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/google/uuid"
)

func (suite *YadosTestSuite) TestPLogIterator() {

	numOfEntries := 100
	ids := make([]string, 0)
	for i := 0; i < numOfEntries; i++ {
		id := uuid.New().String()
		err := suite.cluster.members[0].PLog().Append(&pb.LogEntry{
			Term:    0,
			Index:   0,
			Command: nil,
			Id:      id,
			CmdType: 0,
		})
		if err != nil {
			suite.T().Error(err)
		}
		ids = append(ids, id)
	}
	iter, err := suite.cluster.members[0].PLog().Iterator()
	if err != nil {
		suite.T().Error(err)
	}

	defer func(iter server.PLogIterator) {
		err := iter.Close()
		if err != nil {
			suite.T().Error(err)
		}
	}(iter)

	entry, err := iter.Next()
	if err != nil {
		return
	}

	index := 0
	for entry != nil {
		if entry.Id != ids[index] {
			suite.T().Errorf("Entries mismatch Expected Id %s, Actual Id = %s", ids[index], entry.Id)
		}

		entry, err = iter.Next()
		if err != nil {
			return
		}
		index++
	}
	if index != numOfEntries {
		suite.T().Errorf("Number of Entries mismatch Expected %d, Actual %d", index, numOfEntries)
	}

}
