package events

import pb "github.com/davinash/yados/internal/proto/gen"

//Events struct for events to be used by tests only
type Events struct {
	CommitEntryChan   chan interface{}
	LeaderChangeChan  chan interface{}
	EventSubscription chan interface{}
	PersistEntryChan  chan *pb.WalEntry
}

//NewEvents creates new object of Event interface
func NewEvents() *Events {
	e := &Events{}
	return e
}
