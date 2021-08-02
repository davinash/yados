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
	Peers() []*pb.Peer
	AddPeer(peer *pb.Peer) error
	State() RaftState
	RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error)
	Stop()
	AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error)
	Log() []*pb.LogEntry
}

type raft struct {
	mutex  sync.Mutex
	server Server
	//TODO : Check if we need this or can work with from the object in Server interface
	peers []*pb.Peer

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
	//commitChan         chan<- CommitEntry
}

func (r *raft) Log() []*pb.LogEntry {
	return r.log
}

func (r *raft) Stop() {
	r.Server().Logger().Debug("Stopping Raft Instance")
	if r.state == Dead {
		return
	}
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
		server:        srv,
		quit:          make(chan interface{}),
		triggerAEChan: make(chan struct{}, 1),
	}
	r.peers = make([]*pb.Peer, 0)
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
	for _, peer := range r.Peers() {
		r.nextIndex[peer.Name] = int64(len(r.Log()))
		r.matchIndex[peer.Name] = -1
	}
	r.Server().Logger().Debugf("becomes Leader; term=%d, nextIndex=%v, matchIndex=%v; log=%v", r.currentTerm,
		r.nextIndex, r.matchIndex, r.log)

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
				r.server.Logger().Debug("Stopping startLeader")
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

func (r *raft) leaderSendAEs() {
	r.mutex.Lock()
	savedCurrentTerm := r.currentTerm
	r.mutex.Unlock()

	for _, peer := range r.Peers() {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()

			r.mutex.Lock()
			ni := r.nextIndex[peer.Name]

			var prevLogTerm, prevLogIndex int64
			prevLogIndex = ni - 1

			prevLogTerm = -1
			if prevLogIndex >= 0 {
				prevLogTerm = r.log[prevLogIndex].Term
			}
			entries := r.log[ni:]

			request := pb.AppendEntryRequest{
				Term:         savedCurrentTerm,
				LeaderName:   r.Server().Name(),
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevLogTerm,
				Entries:      entries,
				LeaderCommit: r.commitIndex,
			}
			r.mutex.Unlock()
			r.Server().Logger().Debugf("sending AppendEntries to %s: ni=%d, args=%+v", peer.Name, ni, &request)

			resp, err := r.Server().Send(peer, "RPC.AppendEntries", &request)
			if err != nil {
				r.Server().Logger().Errorf("Sending AppendEntries failed, Error = %v", err)
				return
			}
			r.mutex.Lock()
			defer r.mutex.Unlock()
			reply := resp.(*pb.AppendEntryReply)
			if reply.Term > savedCurrentTerm {
				r.Server().Logger().Debug("term out of date in heartbeat reply")
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
					r.Server().Logger().Debugf("AppendEntries reply from %s success: nextIndex := %v, "+
						"matchIndex := %v; commitIndex := %d", peer.Name, r.nextIndex, r.matchIndex, r.commitIndex)
					if r.commitIndex != savedCommitIndex {
						r.Server().Logger().Debugf("leader sets commitIndex := %d", r.commitIndex)
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
					r.Server().Logger().Debugf("AppendEntries reply from %s !success: nextIndex := %d",
						peer.Name, ni-1)
				}
			}
		}(peer)
	}
}

//
//func (r *raft) commitChanSender() {
//	for range r.newCommitReadyChan {
//		// Find which entries we have to apply.
//		r.mutex.Lock()
//		var savedTerm, savedLastApplied int64
//		savedTerm = r.currentTerm
//		savedLastApplied = r.lastApplied
//		var entries []*pb.LogEntry
//		if r.commitIndex > r.lastApplied {
//			entries = r.log[r.lastApplied+1 : r.commitIndex+1]
//			r.lastApplied = r.commitIndex
//		}
//		r.mutex.Unlock()
//		r.Server().Logger().Debugf("commitChanSender entries=%v, savedLastApplied=%d", entries, savedLastApplied)
//
//		for i, entry := range entries {
//			r.commitChan <- CommitEntry{
//				Command: entry.Command,
//				Index:   savedLastApplied + int64(i) + 1,
//				Term:    savedTerm,
//			}
//		}
//	}
//	r.Server().Logger().Debug("commitChanSender done")
//}

func (r *raft) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	reply := &pb.AppendEntryReply{Id: request.Id}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state == Dead {
		return reply, nil
	}
	r.Server().Logger().Debugf("AppendEntries: %+v", request)

	if request.Term > r.currentTerm {
		r.Server().Logger().Debugf("... term out of date in AppendEntries")
		r.becomeFollower(request.Term)
	}

	if request.Term == r.currentTerm {
		if r.state != Follower {
			r.becomeFollower(request.Term)
		}
		r.electionResetEvent = time.Now()

		// Does our log contain an entry at PrevLogIndex whose term matches
		// PrevLogTerm? Note that in the extreme case of PrevLogIndex=-1 this is
		// vacuously true.
		if request.PrevLogIndex == -1 ||
			(request.PrevLogIndex < int64(len(r.log)) && request.PrevLogTerm == r.log[request.PrevLogIndex].Term) {
			reply.Success = true

			// Find an insertion point - where there's a term mismatch between
			// the existing log starting at PrevLogIndex+1 and the new entries sent
			// in the RPC.
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
			// At the end of this loop:
			// - logInsertIndex points at the end of the log, or an index where the
			//   term mismatches with an entry from the leader
			// - newEntriesIndex points at the end of Entries, or an index where the
			//   term mismatches with the corresponding log entry
			if newEntriesIndex < len(request.Entries) {
				r.Server().Logger().Debugf("... inserting entries %v from index %d", request.Entries[newEntriesIndex:], logInsertIndex)
				r.log = append(r.log[:logInsertIndex], request.Entries[newEntriesIndex:]...)
				r.Server().Logger().Debugf("... log is now: %v", r.log)
			}

			// Set commit index.
			if request.LeaderCommit > r.commitIndex {
				r.commitIndex = intMin(request.LeaderCommit, int64(len(r.log)-1))
				r.Server().Logger().Debugf("... setting commitIndex=%d", r.commitIndex)
				r.newCommitReadyChan <- struct{}{}
			}
		} else {
			// No match for PrevLogIndex/PrevLogTerm. Populate
			// ConflictIndex/ConflictTerm to help the leader bring us up to date
			// quickly.
			if request.PrevLogIndex >= int64(len(r.log)) {
				reply.ConflictIndex = int64(len(r.log))
				reply.ConflictTerm = -1
			} else {
				// PrevLogIndex points within our log, but PrevLogTerm doesn't match
				// r.log[PrevLogIndex].
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
	r.Server().Logger().Debugf("AppendEntries reply: %+v", reply.Term)
	return reply, nil
}

func (r *raft) becomeFollower(term int64) {
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
