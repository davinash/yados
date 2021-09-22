package tests

func (suite *YadosTestSuite) TestWALAppend() {
	storeName := "TestWALAppend"

	if err := CreateSQLStoreAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := CreateTableAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}
	if err := InsertRowsAndWait(suite.cluster, storeName); err != nil {
		suite.T().Fatal(err)
	}

	for _, member := range suite.cluster.members {
		iterator, err := member.WAL().Iterator()
		if err != nil {
			suite.T().Fatal(err)
		}
		entry, err1 := iterator.Next()
		if err1 != nil {
			suite.T().Fatal(err1)
		}
		count := 0
		for entry != nil {
			//printEntry(entry, member.Name(), suite.T())
			count++
			entry, err = iterator.Next()
			if err != nil {
				suite.T().Fatal(err)
			}
		}
		if count != 7 {
			suite.T().Fatalf("[%s] Expected entries = %d Actual = %d", member.Name(), 7, count)
		}

		err = iterator.Close()
		if err != nil {
			suite.T().Fatal(err)
		}
	}
}
