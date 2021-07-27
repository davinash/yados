package server

import (
	"sync"
	"time"
)

//RaftState state of the raft instance
type RaftState int

const (
	//Follower indicate follower state
	Follower RaftState = iota
	//Candidate indicate candidate state
	Candidate
	//Leader indicate leader state
	Leader
	//Dead indicate dead state
	Dead
)

//Raft raft interface
type Raft interface {
	Server() Server
	Peers() []Server
}

type raft struct {
	mutex  sync.Mutex
	server Server
	peers  []Server
	state  RaftState

	votedFor           int
	electionResetEvent time.Time
}

//NewRaft creates new instance of Raft
func NewRaft(srv Server, ready <-chan interface{}) Raft {
	r := &raft{
		server: srv,
	}
	r.peers = make([]Server, 0)
	r.state = Follower
	r.votedFor = -1

	go func() {
		//<-ready
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
	}()

	return r
}

func (r *raft) runElectionTimer() {
}

func (r *raft) Server() Server {
	return r.server
}

func (r *raft) Peers() []Server {
	return r.peers
}
