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
	//EntryPersistEvents events when the entry is persisted
	EntryPersistEvents
)

//Events event interface only to be used by testing module
type Events interface {
	SendEvent(interface{})

	CommitEntryEvent() chan *pb.LogEntry
	LeaderChangeEvent() chan Server
	PersistEntryEvent() chan int

	SetPersistEntryEventThreshold(int)
	PersistEntryEventThreshold() int

	Subscribe(EventType)
	UnSubscribe(EventType)
}

type events struct {
	mutex             sync.Mutex
	commitEntryChan   chan *pb.LogEntry
	leaderChangeChan  chan Server
	eventSubscription map[EventType]bool
	persistEntryChan  chan int
	perstEvThreshold  int
}

//NewEvents creates new object of Event interface
func NewEvents() Events {
	e := &events{
		eventSubscription: make(map[EventType]bool),
		commitEntryChan:   make(chan *pb.LogEntry),
		leaderChangeChan:  make(chan Server),
		persistEntryChan:  make(chan int),
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

func (ev *events) LeaderChangeEvent() chan Server {
	return ev.leaderChangeChan
}

func (ev *events) PersistEntryEvent() chan int {
	return ev.persistEntryChan
}

func (ev *events) SetPersistEntryEventThreshold(count int) {
	ev.perstEvThreshold = count
}

func (ev *events) PersistEntryEventThreshold() int {
	return ev.perstEvThreshold
}

func (ev *events) SendEvent(obj interface{}) {
	ev.mutex.Lock()
	defer ev.mutex.Unlock()

	switch event := obj.(type) {
	case Server:
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
	case int:
		if v, ok := ev.eventSubscription[EntryPersistEvents]; ok {
			if v {
				ev.persistEntryChan <- event
			}
		}
	}
}
