package tests

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"

	"github.com/davinash/yados/internal/server"
)

//WaitForEvents helper function to wait for event on all the members of the cluster
func WaitForEvents(members []server.Server, wg *sync.WaitGroup, evCount int) {
	for _, s := range members {
		s.EventHandler().PersistEntryChan = make(chan *pb.LogEntry)
	}
	for _, member := range members {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			count := 0
			for {
				select {
				case <-s.EventHandler().PersistEntryChan:
					count++
					if count == evCount {
						return
					}
				default:

				}
			}
		}(member)
	}
}

//StopWaitForEvents helper function to stop the wait on channel
func StopWaitForEvents(members []server.Server) {
	for _, s := range members {
		close(s.EventHandler().PersistEntryChan)
		s.EventHandler().PersistEntryChan = nil
	}
}
