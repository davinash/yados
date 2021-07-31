package server

import (
	"context"
	"math/rand"
	"sync"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// CommitEntry is the data reported by Raft to the commit channel. Each commit
// entry notifies the client that consensus was reached on a command and it can
// be applied to the client's state machine.
type CommitEntry struct {
	// Command is the client command being committed.
	Command interface{}
	// Index is the log index at which the client command is committed.
	Index int
	// Term is the Raft term at which the client command is committed.
	Term uint64
}

//LogEntry represents the log entry
type LogEntry struct {
	Command interface{}
	Term    uint64
	Index   uint64
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
	Peers() []*pb.Peer
	AddPeer(peer *pb.Peer) error
	State() RaftState
	RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error)
	Stop()
	AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error)
	Log() []LogEntry
}

type raft struct {
	mutex  sync.Mutex
	server Server
	//TODO : Check if we need this or can work with from the object in Server interface
	peers []*pb.Peer

	// Persistent Raft state on all servers
	currentTerm uint64
	log         []LogEntry
	votedFor    string

	// Volatile Raft state on all servers
	commitIndex        int
	lastApplied        int
	state              RaftState
	electionResetEvent time.Time

	// Volatile Raft state on leaders
	nextIndex  map[int]int
	matchIndex map[int]int

	quit chan interface{}
	wg   sync.WaitGroup
}

func (r *raft) Log() []LogEntry {
	return r.log
}

func (r *raft) Stop() {
	r.Server().Logger().Debug("Stopping Raft Instance")

	close(r.quit)
	r.wg.Wait()

	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.state = Dead

	r.Server().Logger().Debug("Stopped Raft Instance")
}

//NewRaft creates new instance of Raft
func NewRaft(srv Server, peers []*pb.Peer, ready <-chan interface{}) (Raft, error) {
	r := &raft{
		server: srv,
		quit:   make(chan interface{}),
	}
	r.peers = make([]*pb.Peer, 0)
	r.state = Follower
	r.votedFor = ""
	r.commitIndex = -1
	r.lastApplied = -1
	r.nextIndex = make(map[int]int)
	r.matchIndex = make(map[int]int)

	for _, p := range peers {
		_, err := srv.Send(p, "server.AddNewMember", &pb.NewPeerRequest{
			NewPeer: &pb.Peer{
				Name:    srv.Name(),
				Address: srv.Address(),
				Port:    srv.Port(),
			},
		})
		if err != nil {
			return nil, err
		}
		err = r.AddPeer(p)
		if err != nil {
			return nil, err
		}
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		<-ready
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
		r.Server().Logger().Debug("runElectionTimer::NewRaft")
	}()

	return r, nil
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

func (r *raft) State() RaftState {
	return r.state
}

func (r *raft) Server() Server {
	return r.server
}

func (r *raft) Peers() []*pb.Peer {
	return r.peers
}

func (r *raft) AddPeer(newPeer *pb.Peer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.peers = append(r.peers, newPeer)
	return nil
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
		select {
		case <-r.quit:
			r.server.Logger().Debug("Stopping runElectionTimer")
			return
		case <-ticker.C:
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
}

func (r *raft) startElection() {
	r.state = Candidate
	r.currentTerm++
	savedCurrentTerm := r.currentTerm
	r.electionResetEvent = time.Now()
	r.votedFor = r.Server().Name()
	r.server.Logger().Debugf("becomes Candidate (currentTerm=%d); log=%v", savedCurrentTerm, r.log)

	votesReceived := 1
	for _, peer := range r.peers {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()
			r.mutex.Lock()
			savedLastLogIndex, savedLastLogTerm := r.lastLogIndexAndTerm()
			r.mutex.Unlock()

			args := pb.VoteRequest{
				Term:          savedCurrentTerm,
				CandidateName: r.Server().Name(),
				LastLogIndex:  savedLastLogIndex,
				LastLogTerm:   savedLastLogTerm,
			}
			resp, err := r.Server().Send(peer, "RPC.RequestVote", &args)
			if err != nil {
				return
			}
			reply := resp.(*pb.VoteReply)
			r.mutex.Lock()
			defer r.mutex.Unlock()
			r.server.Logger().Debugf("[%s] received RequestVoteReply (%+v, %v)", reply.Id, reply.Term, reply.VoteGranted)
			if r.state != Candidate {
				r.server.Logger().Debugf("[%s] while waiting for reply, state = %v", reply.Id, r.state)
				return
			}
			if reply.Term > savedCurrentTerm {
				r.server.Logger().Debugf("[%s] term out of date in RequestVoteReply", reply.Id)
				r.becomeFollower(reply.Term)
				return
			} else if reply.Term == savedCurrentTerm {
				if reply.VoteGranted {
					votesReceived++
					if votesReceived*2 > len(r.peers)+1 {
						r.server.Logger().Debugf("wins election with %d votes", votesReceived)
						r.startLeader()
						return
					}
				}
			}
		}(peer)
	}
	// Run another election timer, in case this election is not successful.
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
		r.Server().Logger().Debug("runElectionTimer::StartElection")
	}()

}

func (r *raft) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)

	EmptyVoteReply := &pb.VoteReply{Id: request.Id}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state == Dead {
		return EmptyVoteReply, nil
	}
	lastLogIndex, lastLogTerm := r.lastLogIndexAndTerm()
	r.Server().Logger().Debugf("[%s] Received RequestVote: [currentTerm=%d, votedFor=%s] log index/term=(%v, %v)",
		request.Id, r.currentTerm, r.votedFor, lastLogIndex, lastLogTerm)

	if request.Term > r.currentTerm {
		r.Server().Logger().Debugf("[%s] term out of date in RequestVote", request.Id)
		r.becomeFollower(request.Term)
	}

	reply := &pb.VoteReply{Id: request.Id}

	if r.currentTerm == request.Term &&
		(r.votedFor == "" || r.votedFor == request.CandidateName) &&
		(request.LastLogTerm > lastLogTerm ||
			(request.LastLogTerm == lastLogTerm && request.LastLogIndex >= lastLogIndex)) {
		reply.VoteGranted = true
		r.votedFor = request.CandidateName
		r.electionResetEvent = time.Now()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = r.currentTerm
	r.persistToStorage()
	r.Server().Logger().Debugf("[%s] RequestVote reply: [Term = %d, votedFor = %s ]",
		request.Id, reply.Term, r.votedFor)

	return reply, nil
}

func (r *raft) startLeader() {
	r.state = Leader
	r.Server().Logger().Debugf("becomes Leader; term=%d, log=%v", r.currentTerm, r.log)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		// Send periodic heartbeats, as long as still leader.
		for {
			r.leaderSendHeartbeats()
			select {
			case <-r.quit:
				r.server.Logger().Debug("Stopping startLeader")
				return
			case <-ticker.C:
				r.mutex.Lock()
				if r.state != Leader {
					r.mutex.Unlock()
					return
				}
				r.mutex.Unlock()
			}
		}
	}()
}

func (r *raft) leaderSendHeartbeats() {
	r.mutex.Lock()
	savedCurrentTerm := r.currentTerm
	r.mutex.Unlock()

	for _, peer := range r.peers {
		args := pb.AppendEntryRequest{
			Term:       savedCurrentTerm,
			LeaderName: r.Server().Name(),
		}
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()
			resp, err := r.Server().Send(peer, "RPC.AppendEntries", &args)
			if err != nil {
				return
			}
			reply := resp.(*pb.AppendEntryReply)
			r.mutex.Lock()
			defer r.mutex.Unlock()
			if reply.Term > savedCurrentTerm {
				r.Server().Logger().Debug("term out of date in heartbeat reply")
				r.becomeFollower(reply.Term)
				return
			}
		}(peer)
	}
}

func (r *raft) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	EmptyAppendEntryReply := &pb.AppendEntryReply{Id: request.Id}

	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state == Dead {
		return EmptyAppendEntryReply, nil
	}
	r.Server().Logger().Debugf("[%s] AppendEntries: (Term = %v Leader Name = %s)",
		request.Id, request.Term, request.LeaderName)

	if request.Term > r.currentTerm {
		r.Server().Logger().Debugf("[%s]... term out of date in AppendEntries", request.Id)
		r.becomeFollower(request.Term)
	}
	reply := &pb.AppendEntryReply{}

	reply.Success = false
	if request.Term == r.currentTerm {
		if r.state != Follower {
			r.becomeFollower(request.Term)
		}
		r.electionResetEvent = time.Now()
		reply.Success = true
	}

	reply.Term = r.currentTerm
	r.Server().Logger().Debugf("[%s] AppendEntries reply: (Term = %v, Success = %v)", request.Id,
		reply.Term, reply.Success)
	return reply, nil
}

func (r *raft) becomeFollower(term uint64) {
	r.Server().Logger().Debugf("becomes Follower with term=%d; log=%v", term, r.log)
	r.state = Follower
	r.currentTerm = term
	r.votedFor = ""
	r.electionResetEvent = time.Now()

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
		r.Server().Logger().Debug("runElectionTimer::becomeFollower")
	}()

}

func (r *raft) lastLogIndexAndTerm() (int64, int64) {
	if len(r.log) > 0 {
		lastIndex := len(r.log) - 1
		return int64(lastIndex), int64(r.log[lastIndex].Term)
	}
	return -1, -1
}

func (r *raft) persistToStorage() {

}
