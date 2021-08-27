package tests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/davinash/yados/internal/events"

	"github.com/davinash/yados/internal/server"

	"github.com/davinash/yados/cmd/cli/commands/store"
)

func benchmarkPut(numOfPuts int, b *testing.B, cluster *TestCluster, storeName string) {
	for _, s := range cluster.members {
		s.EventHandler().SetPersistEntryEventThreshold(numOfPuts)
		s.EventHandler().Subscribe(events.EntryPersistEvents)
	}
	defer func() {
		for _, s := range cluster.members {
			s.EventHandler().UnSubscribe(events.EntryPersistEvents)
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
		key := fmt.Sprintf("Key-%d", i)
		val := fmt.Sprintf("Value-%d", i)

		err := server.ExecuteCmdPut(&server.PutArgs{
			Address:   cluster.members[0].Address(),
			Port:      cluster.members[0].Port(),
			Key:       key,
			Value:     val,
			StoreName: storeName,
		})
		if err != nil {
			b.Fatal(err)
		}
		b.Logf("Put Success  Key = %s Value = %s", key, val)
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

	b.Run("ABC", func(b *testing.B) {
		storeName := "BenchmarkPut"
		b.Log("Creating a store")
		err = server.ExecuteCmdCreateStore(&server.CreateCommandArgs{
			Address: cluster.members[0].Address(),
			Port:    cluster.members[0].Port(),
			Name:    storeName,
		})
		if err != nil {
			b.Fatal(err)
		}

		benchmarkPut(100000, b, cluster, storeName)

		b.Log("Deleting a store")

		err := store.ExecuteCmdDeleteStore(&store.DeleteArgs{
			Address:   cluster.members[0].Address(),
			Port:      cluster.members[0].Port(),
			StoreName: storeName,
		})
		if err != nil {
			b.Fatal(err)
		}
	})
}
