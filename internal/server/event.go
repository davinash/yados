package server

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//EventType event type
type EventType int

const (
	//CommitEntryEvents type for commit entry
	CommitEntryEvents = iota
	//LeaderChangeEvents type when the leader ship changes
	LeaderChangeEvents
)

//Events event interface only to be used by testing module
type Events interface {
	SendEvent(interface{})

	CommitEntryEvent() chan *pb.LogEntry
	LeaderChangeEvent() chan *pb.Peer

	Subscribe(EventType)
	UnSubscribe(EventType)
}

type events struct {
	mutex             sync.Mutex
	commitEntryChan   chan *pb.LogEntry
	leaderChangeChan  chan *pb.Peer
	eventSubscription map[EventType]bool
}

//NewEvents creates new object of Event interface
func NewEvents() Events {
	e := &events{
		eventSubscription: make(map[EventType]bool),
		commitEntryChan:   make(chan *pb.LogEntry),
		leaderChangeChan:  make(chan *pb.Peer),
	}
	return e
}

func (ev *events) Subscribe(evType EventType) {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()
	ev.eventSubscription[evType] = true
}

func (ev *events) UnSubscribe(evType EventType) {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()
	ev.eventSubscription[evType] = false
}

func (ev *events) CommitEntryEvent() chan *pb.LogEntry {
	return ev.commitEntryChan
}

func (ev *events) LeaderChangeEvent() chan *pb.Peer {
	return ev.leaderChangeChan
}

func (ev *events) SendEvent(obj interface{}) {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()

	switch event := obj.(type) {
	case *pb.Peer:
		if v, ok := ev.eventSubscription[LeaderChangeEvents]; ok {
			if v {
				ev.leaderChangeChan <- event
			}
		}
	case *pb.LogEntry:
		if v, ok := ev.eventSubscription[CommitEntryEvents]; ok {
			if v {
				ev.commitEntryChan <- event
			}
		}
	}
}
