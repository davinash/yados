package server

import (
	"math/rand"
	"sync"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//LogEntry represents the log entry
type LogEntry struct {
	Command interface{}
	Term    int
}

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
	//TODO : Check if we need this or can work with from the object in Server interface
	peers []Server
	state RaftState

	currentTerm uint64
	log         []LogEntry

	votedFor           string
	electionResetEvent time.Time
}

//NewRaft creates new instance of Raft
func NewRaft(srv Server, ready <-chan interface{}) Raft {
	r := &raft{
		server: srv,
	}
	r.peers = make([]Server, 0)
	r.state = Follower
	r.votedFor = ""

	go func() {
		<-ready
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
	}()

	return r
}

func (s RaftState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Dead:
		return "Dead"
	default:
		panic("unreachable")
	}
}

func (r *raft) electionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (r *raft) runElectionTimer() {
	timeoutDuration := r.electionTimeout()
	r.mutex.Lock()
	termStarted := r.currentTerm
	r.mutex.Unlock()
	r.server.Logger().Debugf("election timer started (%v), term=%d", timeoutDuration, termStarted)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		<-ticker.C
		r.mutex.Lock()
		if r.state != Candidate && r.state != Follower {
			r.server.Logger().Debugf("in election timer state=%s, bailing out", r.state)
			r.mutex.Unlock()
			return
		}
		if termStarted != r.currentTerm {
			r.server.Logger().Debugf("in election timer term changed from %d to %d, bailing out",
				termStarted, r.currentTerm)
			r.mutex.Unlock()
			return
		}
		if elapsed := time.Since(r.electionResetEvent); elapsed >= timeoutDuration {
			r.startElection()
			r.mutex.Unlock()
			return
		}

		r.mutex.Unlock()
	}
}

func (r *raft) startElection() {
	r.state = Candidate
	r.currentTerm++
	savedCurrentTerm := r.currentTerm
	r.electionResetEvent = time.Now()
	r.votedFor = r.server.Name()
	r.server.Logger().Debugf("becomes Candidate (currentTerm=%d); log=%v", savedCurrentTerm, r.log)
	//votesReceived := 1
	for _, peer := range r.peers {
		go func(peer Server) {
			args := pb.VoteRequest{
				Term:          savedCurrentTerm,
				CandidateName: r.Server().Name(),
			}
			//var reply *pb.VoteReply
			if r, err := r.Server().Send(peer, "RPC.RequestVote", &args); err == nil {
				_ = r.(*pb.VoteReply)
			}
		}(peer)
	}

}

func (r *raft) Server() Server {
	return r.server
}

func (r *raft) Peers() []Server {
	return r.peers
}
