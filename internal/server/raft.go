package server

import (
	"context"
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
	Peers() []*pb.Peer
	AddPeer(peer *pb.Peer) error
	State() RaftState
	RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error)
	Stop()
	AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error)
}

type raft struct {
	mutex  sync.Mutex
	server Server
	//TODO : Check if we need this or can work with from the object in Server interface
	peers []*pb.Peer
	state RaftState

	currentTerm uint64
	log         []LogEntry

	votedFor           string
	electionResetEvent time.Time
}

func (r *raft) Stop() {
	r.state = Dead
}

//NewRaft creates new instance of Raft
func NewRaft(srv Server, peers []*pb.Peer, ready <-chan interface{}) (Raft, error) {
	r := &raft{
		server: srv,
	}
	r.peers = make([]*pb.Peer, 0)
	r.state = Follower
	r.votedFor = ""

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

	go func() {
		<-ready
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
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
	r.votedFor = r.Server().Name()
	r.server.Logger().Debugf("becomes Candidate (currentTerm=%d); log=%v", savedCurrentTerm, r.log)

	votesReceived := 1
	for _, peer := range r.peers {
		go func(peer *pb.Peer) {
			args := pb.VoteRequest{
				Term:          savedCurrentTerm,
				CandidateName: r.Server().Name(),
			}
			resp, err := r.Server().Send(peer, "RPC.RequestVote", &args)
			if err != nil {
				return
			}
			reply := resp.(*pb.VoteReply)
			r.mutex.Lock()
			defer r.mutex.Unlock()
			r.server.Logger().Debugf("received RequestVoteReply %+v", reply)
			if r.state != Candidate {
				r.server.Logger().Debugf("while waiting for reply, state = %v", r.state)
				return
			}
			if reply.Term > savedCurrentTerm {
				r.server.Logger().Debug("term out of date in RequestVoteReply")
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
	go r.runElectionTimer()
}

func (r *raft) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)

	EmptyVoteReply := &pb.VoteReply{Id: request.Id}

	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state == Dead {
		return EmptyVoteReply, nil
	}
	r.Server().Logger().Debugf("[%s] Received RequestVote: [currentTerm=%d, votedFor=%s]",
		request.Id, r.currentTerm, r.votedFor)

	if request.Term > r.currentTerm {
		r.Server().Logger().Debugf("[%s] term out of date in RequestVote", request.Id)
		r.becomeFollower(request.Term)
	}
	reply := &pb.VoteReply{Id: request.Id}
	if r.currentTerm == request.Term && (r.votedFor == "" || r.votedFor == request.CandidateName) {
		reply.VoteGranted = true
		r.votedFor = request.CandidateName
		r.electionResetEvent = time.Now()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = r.currentTerm
	r.Server().Logger().Debugf("[%s] RequestVote reply: [Term = %d, votedFor = %s ]",
		request.Id, reply.Term, r.votedFor)

	return reply, nil
}

func (r *raft) startLeader() {
	r.state = Leader
	r.Server().Logger().Debugf("becomes Leader; term=%d, log=%v", r.currentTerm, r.log)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		// Send periodic heartbeats, as long as still leader.
		for {
			r.leaderSendHeartbeats()
			<-ticker.C

			r.mutex.Lock()
			if r.state != Leader {
				r.mutex.Unlock()
				return
			}
			r.mutex.Unlock()
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
		go func(peer *pb.Peer) {
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

	go r.runElectionTimer()
}
