package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/davinash/yados/internal/server"

	"github.com/davinash/yados/cmd/cli/commands/store"
)

func benchmarkPut(numOfPuts int, b *testing.B, cluster *TestCluster, storeName string) {
	for _, s := range cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(server.EntryPersistEvents)
	}
	defer func() {
		for _, s := range cluster.members {
			s.EventHandler().UnSubscribe(server.EntryPersistEvents)
		}
	}()

	wg := sync.WaitGroup{}
	for _, member := range cluster.members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			<-s.EventHandler().PersistEntryEvent()
		}(member)
	}

	for i := 0; i < numOfPuts; i++ {
		err := store.ExecutePutCommand(&store.PutArgs{
			Address:   cluster.members[0].Address(),
			Port:      cluster.members[0].Port(),
			Key:       fmt.Sprintf("Key-%d", i),
			Value:     fmt.Sprintf("Value-%d", i),
			StoreName: storeName,
		})
		if err != nil {
			b.Error(err)
		}
	}
	// Wait for replication to happen
	wg.Wait()
}

func BenchmarkPut(b *testing.B) {
	cluster := &TestCluster{}
	logDir := SetupDataDirectory()

	err := CreateNewClusterEx(3, cluster, logDir, "info")
	if err != nil {
		b.Fail()
	}
	WaitForLeaderElection(cluster)

	defer func() {
		StopCluster(cluster)
		Cleanup(logDir)
	}()

	storeName := "BenchmarkPut"
	err = store.CreateCommandExecute(&store.CreateCommandArgs{
		Address: cluster.members[0].Address(),
		Port:    cluster.members[0].Port(),
		Name:    storeName,
	})
	if err != nil {
		b.Error(err)
	}

	b.Run("ABC", func(b *testing.B) {
		benchmarkPut(100, b, cluster, storeName)
	})
}
