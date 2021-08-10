package server

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/davinash/yados/internal/proto/gen"
)

// CommitEntry is the data reported by Raft to the commit channel. Each commit
// entry notifies the client that consensus was reached on a command and it can
// be applied to the client's state machine.
type CommitEntry struct {
	// Command is the client command being committed.
	Command interface{}
	// Index is the log index at which the client command is committed.
	Index int64
	// Term is the Raft term at which the client command is committed.
	Term int64
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
	Peers() map[string]*pb.Peer
	AddPeer(peer *pb.Peer) error
	State() RaftState
	RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error)
	RemovePeer(*pb.RemovePeerRequest) error
	Start()
	Stop()
	AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error)
	Log() []*pb.LogEntry
	Submit([]byte) error
}

type raft struct {
	mutex  sync.Mutex
	server Server
	//TODO : Check if we need this or can work with from the object in Server interface
	peers map[string]*pb.Peer

	// Persistent Raft state on all servers
	currentTerm int64
	log         []*pb.LogEntry
	votedFor    string

	// Volatile Raft state on all servers
	commitIndex        int64
	lastApplied        int64
	state              RaftState
	electionResetEvent time.Time

	// Volatile Raft state on leaders
	nextIndex  map[string]int64
	matchIndex map[string]int64

	triggerAEChan chan struct{}

	quit chan interface{}
	wg   sync.WaitGroup

	newCommitReadyChan chan struct{}
	logger             *logrus.Entry
}

func (r *raft) Log() []*pb.LogEntry {
	return r.log
}

func (r *raft) Stop() {
	r.logger.Debug("Stopping Raft Instance")
	if r.state == Dead {
		return
	}
	close(r.quit)

	err := r.RemoveSelf()
	if err != nil {
		return
	}
	r.wg.Wait()
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.state = Dead

	r.logger.Debug("Stopped Raft Instance")
}

//NewRaft creates new instance of Raft
func NewRaft(srv Server, peers []*pb.Peer) (Raft, error) {
	r := &raft{
		server:        srv,
		quit:          make(chan interface{}),
		triggerAEChan: make(chan struct{}, 1),
		logger:        srv.Logger(),
	}
	r.peers = make(map[string]*pb.Peer)
	r.state = Follower
	r.votedFor = ""
	r.commitIndex = -1
	r.lastApplied = -1
	r.nextIndex = make(map[string]int64)
	r.matchIndex = make(map[string]int64)
	r.newCommitReadyChan = make(chan struct{}, 16)

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
		// add this peer to self
		err = r.AddPeer(p)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r *raft) Start() {
	go func() {
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
		r.logger.Debug("runElectionTimer::NewRaft")
	}()
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

func (r *raft) Peers() map[string]*pb.Peer {
	return r.peers
}

func (r *raft) AddPeer(newPeer *pb.Peer) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.peers[newPeer.Name] = newPeer
	return nil
}

func (r *raft) RemovePeer(request *pb.RemovePeerRequest) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.logger.Debugf("[%s] RemovePeer: Received Name = %s; Address = %s; Port = %d;",
		request.Id, request.GetPeer().Name, request.GetPeer().Address, request.GetPeer().Port)

	delete(r.peers, request.GetPeer().Name)

	return nil
}

func (r *raft) RemoveSelf() error {
	for _, peer := range r.peers {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()
			args := pb.RemovePeerRequest{
				Peer: r.Server().Self(),
			}

			resp, err := r.Server().Send(peer, "RPC.RemoveSelf", &args)
			if err != nil {
				r.logger.Errorf("failed to send RemoveSelf to %s, Error = %v", peer.Name, err)
				return
			}
			_ = resp.(*pb.RemovePeerReply)

		}(peer)
	}
	return nil
}

func (r *raft) electionTimeout() time.Duration {
	return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

//ErrorNotALeader error returned where no server leader
var ErrorNotALeader = errors.New("now a leader")

func (r *raft) Submit(command []byte) error {
	r.logger.Debug("Entering Submit")
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != Leader {
		return ErrorNotALeader
	}
	r.log = append(r.log, &pb.LogEntry{Command: command, Term: r.currentTerm})
	r.persistToStorage()
	r.triggerAEChan <- struct{}{}
	r.logger.Debug("Exiting Submit")
	return nil
}

func (r *raft) runElectionTimer() {
	timeoutDuration := r.electionTimeout()
	r.mutex.Lock()
	termStarted := r.currentTerm
	r.mutex.Unlock()
	r.logger.Debugf("election timer started (%v), term=%d", timeoutDuration, termStarted)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-r.quit:
			r.logger.Debug("Stopping runElectionTimer")
			return
		case <-ticker.C:
			r.startOrIgnoreElection(termStarted, timeoutDuration)
		}
	}
}

func (r *raft) startOrIgnoreElection(termStarted int64, timeoutDuration time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != Candidate && r.state != Follower {
		r.logger.Tracef("in election timer state=%s, bailing out", r.state)
		return
	}

	if termStarted != r.currentTerm {
		r.logger.Tracef("in election timer term changed from %d to %d, bailing out",
			termStarted, r.currentTerm)
		return
	}

	if elapsed := time.Since(r.electionResetEvent); elapsed >= timeoutDuration {
		r.startElection()
		return
	}
}

func (r *raft) processVotingReply(reply *pb.VoteReply, votesReceived *int, savedCurrentTerm *int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != Candidate {
		r.logger.Debugf("[%s] while waiting for reply, state = %v", reply.Id, r.state)
		return
	}

	if reply.Term > *savedCurrentTerm {
		r.logger.Debugf("[%s] term out of date in RequestVoteReply", reply.Id)
		r.becomeFollower(reply.Term)
		return
	} else if reply.Term == *savedCurrentTerm {
		if reply.VoteGranted {
			*votesReceived++
			if *votesReceived*2 > len(r.peers)+1 {
				r.logger.Infof("wins election with %d votes", *votesReceived)
				r.startLeader()
				return
			}
		}
	}
}

func (r *raft) startElection() {
	// Make yourself candidate
	r.state = Candidate
	// Increase the current term
	r.currentTerm++
	savedCurrentTerm := r.currentTerm
	r.electionResetEvent = time.Now()
	// vote for yourself
	r.votedFor = r.Server().Name()
	r.logger.Tracef("becomes Candidate (currentTerm=%d); log=%v", savedCurrentTerm, r.log)

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
			// Send request vote to all peer
			resp, err := r.Server().Send(peer, "RPC.RequestVote", &args)
			if err != nil {
				r.logger.Errorf("failed to send RequestVote to %s, Error = %v", peer.Name, err)
				return
			}
			reply := resp.(*pb.VoteReply)
			r.logger.Debugf("[%s] received RequestVoteReply (%+v, %v)", reply.Id, reply.Term, reply.VoteGranted)

			r.processVotingReply(reply, &votesReceived, &savedCurrentTerm)
		}(peer)
	}
	// Run another election timer, in case this election is not successful.
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
		r.logger.Debug("runElectionTimer::StartElection")
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
	r.logger.Debugf("[%s] RequestVote: Received currentTerm=%d; votedFor=%s; lastLogIndex=%v; lastLogTerm %v",
		request.Id, r.currentTerm, r.votedFor, lastLogIndex, lastLogTerm)

	if request.Term > r.currentTerm {
		r.logger.Debugf("[%s] RequestVote: request.Term =%v; r.currentTerm =%v; "+
			"term out of date in RequestVote", request.Id, request.Term, r.currentTerm)
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
	reply.CandidateName = r.Server().Name()

	r.persistToStorage()

	r.logger.Debugf("[%s] RequestVote: Reply Term = %v; VoteGranted = %v; CandidateName = %s",
		reply.Id, reply.Term, reply.VoteGranted, reply.CandidateName)

	return reply, nil
}

func (r *raft) startLeader() {
	r.state = Leader

	r.Server().SetLeader(r.Server().Self())

	for _, peer := range r.Peers() {
		r.nextIndex[peer.Name] = int64(len(r.Log()))
		r.matchIndex[peer.Name] = -1
	}
	r.logger.Debugf("becomes Leader; term=%d, nextIndex=%v, matchIndex=%v;", r.currentTerm,
		r.nextIndex, r.matchIndex)

	r.wg.Add(1)
	go func(heartbeatTimeout time.Duration) {
		defer r.wg.Done()
		// Immediately send AEs to peers.
		r.leaderSendAEs()
		t := time.NewTimer(heartbeatTimeout)
		defer t.Stop()
		for {
			doSend := false
			select {
			case <-r.quit:
				r.logger.Debug("Stopping startLeader")
				return
			case <-t.C:
				doSend = true
				// Reset timer to fire again after heartbeatTimeout.
				t.Stop()
				t.Reset(heartbeatTimeout)
			case _, ok := <-r.triggerAEChan:
				if ok {
					doSend = true
				} else {
					return
				}
				// Reset timer for heartbeatTimeout.
				if !t.Stop() {
					<-t.C
				}
				t.Reset(heartbeatTimeout)
			}
			if doSend {
				r.mutex.Lock()
				if r.state != Leader {
					r.mutex.Unlock()
					return
				}
				r.mutex.Unlock()
				r.leaderSendAEs()
			}
		}
	}(50 * time.Millisecond)
}

func (r *raft) processAEReply(resp interface{}, savedCurrentTerm int64, peer *pb.Peer, ni int64, entries []*pb.LogEntry) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	reply := resp.(*pb.AppendEntryReply)
	if reply.Term > savedCurrentTerm {
		r.logger.Debugf("[%s] Reply.Term = %v; SavedCurrentTerm = %v; term out of date in heartbeat reply",
			reply.Id, reply.Term, savedCurrentTerm)
		r.becomeFollower(reply.Term)
		return
	}
	if r.state == Leader && savedCurrentTerm == reply.Term {
		if reply.Success {
			r.nextIndex[peer.Name] = ni + int64(len(entries))
			r.matchIndex[peer.Name] = r.nextIndex[peer.Name] - 1
			savedCommitIndex := r.commitIndex
			for i := r.commitIndex + 1; i < int64(len(r.log)); i++ {
				if r.log[i].Term == r.currentTerm {
					matchCount := 1
					for _, p := range r.Peers() {
						if r.matchIndex[p.Name] >= i {
							matchCount++
						}
					}
					if matchCount*2 > len(r.Peers())+1 {
						r.commitIndex = i
					}
				}
			}
			r.logger.Debugf("[%s] AppendEntries reply from %s success: nextIndex := %v, "+
				"matchIndex := %v; commitIndex := %d", reply.Id, peer.Name, r.nextIndex, r.matchIndex, r.commitIndex)
			if r.commitIndex != savedCommitIndex {
				r.logger.Debugf("[%s] leader sets commitIndex := %d", reply.Id, r.commitIndex)
				r.newCommitReadyChan <- struct{}{}
				r.triggerAEChan <- struct{}{}
			}
		} else {
			if reply.ConflictTerm >= 0 {
				var lastIndexOfTerm int64
				lastIndexOfTerm = -1
				for i := len(r.log) - 1; i >= 0; i-- {
					if r.log[i].Term == reply.ConflictTerm {
						lastIndexOfTerm = int64(i)
						break
					}
				}
				if lastIndexOfTerm >= 0 {
					r.nextIndex[peer.Name] = lastIndexOfTerm + 1
				} else {
					r.nextIndex[peer.Name] = reply.ConflictIndex
				}
			} else {
				r.nextIndex[peer.Name] = reply.ConflictIndex
			}
			r.logger.Debugf("[%s] AppendEntries reply from %s !success: nextIndex := %d",
				reply.Id, peer.Name, ni-1)
		}
	}
}

func (r *raft) createAERequest(peer *pb.Peer, savedCurrentTerm int64) (*pb.AppendEntryRequest, int64, []*pb.LogEntry) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	nextIndex := r.nextIndex[peer.Name]
	var prevLogTerm, prevLogIndex int64
	prevLogIndex = nextIndex - 1
	prevLogTerm = -1
	if prevLogIndex >= 0 {
		prevLogTerm = r.log[prevLogIndex].Term
	}
	entries := r.log[nextIndex:]

	request := pb.AppendEntryRequest{
		Term:         savedCurrentTerm,
		Leader:       r.Server().Self(),
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: r.commitIndex,
		NextIndex:    nextIndex,
	}
	return &request, nextIndex, entries
}

func (r *raft) leaderSendAEs() {
	r.mutex.Lock()
	savedCurrentTerm := r.currentTerm
	r.mutex.Unlock()

	for _, peer := range r.Peers() {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()

			request, nextIndex, entries := r.createAERequest(peer, savedCurrentTerm)

			resp, err := r.Server().Send(peer, "RPC.AppendEntries", request)
			if err != nil {
				r.logger.Errorf("[%s] Sending AppendEntries failed, Error = %v", request.Id, err)
				return
			}
			r.processAEReply(resp, savedCurrentTerm, peer, nextIndex, entries)
		}(peer)
	}
}

func (r *raft) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	reply := &pb.AppendEntryReply{Id: request.Id}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state == Dead {
		return reply, nil
	}
	r.Server().SetLeader(request.Leader)

	r.logger.Debugf("[%s] [%s]  Received AppendEntries, current Term = %v", r.state,
		request.Id, r.currentTerm)

	if request.Term > r.currentTerm {
		r.logger.Debugf("[%s] [%s] request.Term = %v; currentTerm = %v; term out of date in AppendEntries",
			r.state, request.Id, request.Term, r.currentTerm)
		r.becomeFollower(request.Term)
	}

	if request.Term == r.currentTerm {
		if r.state != Follower {
			r.becomeFollower(request.Term)
		}
		r.electionResetEvent = time.Now()

		if request.PrevLogIndex == -1 ||
			(request.PrevLogIndex < int64(len(r.log)) && request.PrevLogTerm == r.log[request.PrevLogIndex].Term) {
			reply.Success = true

			logInsertIndex := request.PrevLogIndex + 1
			newEntriesIndex := 0

			for {
				if logInsertIndex >= int64(len(r.log)) || newEntriesIndex >= len(request.Entries) {
					break
				}
				if r.log[logInsertIndex].Term != request.Entries[newEntriesIndex].Term {
					break
				}
				logInsertIndex++
				newEntriesIndex++
			}
			if newEntriesIndex < len(request.Entries) {
				r.logger.Debugf("[%s] [%s] inserting entries %v from index %d", r.state,
					request.Id, request.Entries[newEntriesIndex:], logInsertIndex)
				r.log = append(r.log[:logInsertIndex], request.Entries[newEntriesIndex:]...)
				r.logger.Debugf("[%s] [%s] log is now: %v", r.state, request.Id, r.log)
			}

			if request.LeaderCommit > r.commitIndex {
				r.commitIndex = intMin(request.LeaderCommit, int64(len(r.log)-1))
				r.logger.Debugf("[%s] [%s] setting commitIndex=%d", r.state,
					request.Id, r.commitIndex)
				r.newCommitReadyChan <- struct{}{}
			}
		} else {
			if request.PrevLogIndex >= int64(len(r.log)) {
				reply.ConflictIndex = int64(len(r.log))
				reply.ConflictTerm = -1
			} else {
				reply.ConflictTerm = r.log[request.PrevLogIndex].Term

				var i int64
				for i = request.PrevLogIndex - 1; i >= 0; i-- {
					if r.log[i].Term != reply.ConflictTerm {
						break
					}
				}
				reply.ConflictIndex = i + 1
			}
		}
	}

	reply.Term = r.currentTerm
	r.persistToStorage()
	r.logger.Debugf("[%s] [%s] AppendEntries reply: Term = %v; ConflictTerm =%v; ConflictIndex =%v; Success =%v;",
		r.state, request.Id, reply.Term, reply.ConflictTerm, reply.ConflictIndex, reply.Success)
	return reply, nil
}

func (r *raft) becomeFollower(term int64) {
	r.logger.Debugf("becomes Follower with term=%d; log=%v", term, r.log)
	r.state = Follower
	r.currentTerm = term
	r.votedFor = ""
	r.electionResetEvent = time.Now()

	r.Server().SetLeader(nil)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
		r.logger.Debug("runElectionTimer::becomeFollower")
	}()
}

func (r *raft) lastLogIndexAndTerm() (int64, int64) {
	if len(r.log) > 0 {
		var lastIndex, term int64

		lastIndex = int64(len(r.log) - 1)
		term = r.log[lastIndex].Term

		return lastIndex, term
	}
	return -1, -1
}

func (r *raft) persistToStorage() {
}

func intMin(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
