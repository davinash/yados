package raft

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/davinash/yados/internal/store"

	"github.com/davinash/yados/internal/events"
	"github.com/davinash/yados/internal/rpc"
	"github.com/davinash/yados/internal/wal"

	"github.com/google/uuid"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

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

//State of the raft instance
type State int

const (
	//Follower indicate follower state
	Follower State = iota
	//Candidate indicate candidate state
	Candidate
	//Leader indicate leader state
	Leader
	//Dead indicate dead state
	Dead
)

//Raft raft interface
type Raft interface {
	Peers() map[string]*pb.Peer
	AddPeer(peer *pb.Peer) error
	State() State
	RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error)
	RemovePeer(*pb.RemovePeerRequest) error
	Start()
	Stop()
	AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error)
	Log() []*pb.WalEntry
	SubmitAndWait(interface{}, string, pb.CommandType) error
	AddCommandListener(string) error
	WaitForCommandCompletion(string) error

	Server() *pb.Peer

	EventHandler() *events.Events

	IsTestMode() bool

	StorageManager() store.Manager
}

type raft struct {
	mutex sync.Mutex
	peers map[string]*pb.Peer

	// Persistent Raft state on all servers
	currentTerm int64
	log         []*pb.WalEntry
	votedFor    string

	// Volatile Raft state on all servers
	commitIndex        int64
	lastApplied        int64
	state              State
	electionResetEvent time.Time

	// Volatile Raft state on leaders
	nextIndex  map[string]int64
	matchIndex map[string]int64

	triggerAEChan  chan struct{}
	commitsChanMap map[string]chan error

	quit chan interface{}
	wg   sync.WaitGroup

	newCommitReadyChan chan struct{}
	logger             *logrus.Logger
	wal                wal.Wal
	rpcServer          rpc.Server
	server             *pb.Peer
	ev                 *events.Events
	isTestMode         bool
	storageMgr         store.Manager
}

//Args argument structure for Raft Instance
type Args struct {
	IsTestMode   bool
	Logger       *logrus.Logger
	PstLog       wal.Wal
	RPCServer    rpc.Server
	Server       *pb.Peer
	EventHandler *events.Events
	StorageMgr   store.Manager
}

//NewRaft creates new instance of Raft
func NewRaft(args *Args) (Raft, error) {
	r := &raft{
		quit:          make(chan interface{}),
		triggerAEChan: make(chan struct{}, 1),
		logger:        args.Logger,
		wal:           args.PstLog,
		rpcServer:     args.RPCServer,
		server:        args.Server,
		ev:            args.EventHandler,
		isTestMode:    args.IsTestMode,
		storageMgr:    args.StorageMgr,
	}

	r.peers = make(map[string]*pb.Peer)
	r.state = Follower
	r.commitIndex = -1
	r.lastApplied = -1
	r.nextIndex = make(map[string]int64)
	r.matchIndex = make(map[string]int64)
	r.newCommitReadyChan = make(chan struct{}, 16)
	r.commitsChanMap = make(map[string]chan error)

	err := r.InitializeFromStorage()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *raft) InitializeFromStorage() error {
	// Apply term and voted for state from the
	// persistent log
	term, votedFor, err := r.wal.ReadState()
	if err != nil {
		return err
	}
	r.currentTerm = term
	r.votedFor = votedFor

	iter, err1 := r.wal.Iterator()
	if err1 != nil {
		return nil
	}
	defer func(iter wal.Iterator) {
		err := iter.Close()
		if err != nil {
			r.logger.Warnf("Failed to close the iterator, Error = %v", err)
		}
	}(iter)

	entry, err2 := iter.Next()
	if err2 != nil {
		return nil
	}

	for entry != nil {
		r.log = append(r.log, entry)
		err := r.StorageManager().Apply(entry)
		if err != nil {
			return err
		}
		entry, err = iter.Next()
		if err != nil {
			return nil
		}
	}
	return nil
}

func (r *raft) Start() {
	go func() {
		r.mutex.Lock()
		r.electionResetEvent = time.Now()
		r.mutex.Unlock()
		r.runElectionTimer()
	}()
	go r.commitChanSender()
}

func (r *raft) Stop() {
	r.logger.Debugf("[%s] Stopping Raft Instance", r.server.Name)
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

	r.logger.Debugf("[%s] Stopped Raft Instance", r.server.Name)
}

func (s State) String() string {
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

func (r *raft) Server() *pb.Peer {
	return r.server
}

func (r *raft) Log() []*pb.WalEntry {
	return r.log
}

func (r *raft) State() State {
	return r.state
}

func (r *raft) Peers() map[string]*pb.Peer {
	return r.peers
}

func (r *raft) EventHandler() *events.Events {
	return r.ev
}

func (r *raft) StorageManager() store.Manager {
	return r.storageMgr
}

func (r *raft) IsTestMode() bool {
	return r.isTestMode
}

func (r *raft) AddCommandListener(id string) error {
	r.commitsChanMap[id] = make(chan error)
	return nil
}

func (r *raft) WaitForCommandCompletion(id string) error {
	var err error
	if v, ok := r.commitsChanMap[id]; ok {
		err = <-v
	}
	return err
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
				Peer: r.server,
			}

			resp, err := r.rpcServer.Send(peer, "RPC.RemoveSelf", &args)
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

func (r *raft) runElectionTimer() {
	timeoutDuration := r.electionTimeout()
	r.mutex.Lock()
	termStarted := r.currentTerm
	r.mutex.Unlock()
	r.logger.Debugf("[%s] election timer started (%v), term=%d", r.Server().Name, timeoutDuration, termStarted)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-r.quit:
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
		//r.logger.Debugf("[%s] in election timer term state= %s;", r.Server().Name, r.state.String())
		return
	}

	if termStarted != r.currentTerm {
		//r.logger.Debugf("[%s] in election timer term changed from %d to %d, bailing out",
		//r.Server().Name, termStarted, r.currentTerm)
		return
	}

	if elapsed := time.Since(r.electionResetEvent); elapsed >= timeoutDuration {
		r.logger.Debugf("[%s] Starting new election state=%v; currentTerm=%v;", r.Server().Name,
			r.State(), r.currentTerm)
		r.startElection()
		return
	}
}

func (r *raft) processVotingReply(reply *pb.VoteReply, votesReceived *int, savedCurrentTerm *int64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != Candidate {
		return
	}

	if reply.Term > *savedCurrentTerm {
		r.logger.Debugf("[%s][%s] term out of date in RequestVoteReply", r.Server().Name, reply.Id)
		r.becomeFollower(reply.Term)
		return
	} else if reply.Term == *savedCurrentTerm {
		if reply.VoteGranted {
			*votesReceived++
			if *votesReceived*2 > len(r.peers)+1 {
				r.logger.Infof("[%s] wins election with %d votes", r.server.Name, *votesReceived)
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
	r.votedFor = r.server.Name

	votesReceived := 1
	voteReqId := uuid.New().String()
	for _, peer := range r.peers {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()

			r.mutex.Lock()
			savedLastLogIndex, savedLastLogTerm := r.lastLogIndexAndTerm()
			r.mutex.Unlock()

			args := pb.VoteRequest{
				Term:          savedCurrentTerm,
				CandidateName: r.server.Name,
				LastLogIndex:  savedLastLogIndex,
				LastLogTerm:   savedLastLogTerm,
				Id:            voteReqId,
			}
			// Send request vote to all peer
			resp, err := r.rpcServer.Send(peer, "RPC.RequestVote", &args)
			if err != nil {
				r.logger.Errorf("[%s] failed to send RequestVote to %s, Error = %v", r.server.Name, peer.Name, err)
				return
			}
			reply := resp.(*pb.VoteReply)
			r.logger.Debugf("[%s] [%s] received RequestVoteReply from %v (%v, %v)", r.Server().Name,
				reply.Id, peer.Name, reply.Term, reply.VoteGranted)

			r.processVotingReply(reply, &votesReceived, &savedCurrentTerm)
		}(peer)
	}
	// Run another election timer, in case this election is not successful.
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
	}()

}

func (r *raft) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	//time.Sleep(time.Duration(1+rand.Intn(5)) * time.Millisecond)

	EmptyVoteReply := &pb.VoteReply{Id: request.Id}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state == Dead {
		return EmptyVoteReply, nil
	}
	lastLogIndex, lastLogTerm := r.lastLogIndexAndTerm()
	r.logger.Debugf("[%s] [%s] Received RequestVote from=%s currentTerm=%d; votedFor=%s; lastLogIndex=%v; lastLogTerm %v",
		r.Server().Name, request.Id, request.CandidateName, r.currentTerm, r.votedFor, lastLogIndex, lastLogTerm)

	if request.Term > r.currentTerm {
		r.logger.Debugf("[%s] [%s] RequestVote: request.Term =%v; r.currentTerm =%v; "+
			"term out of date in RequestVote", r.Server().Name, request.Id, request.Term, r.currentTerm)
		r.becomeFollower(request.Term)
	}
	reply := &pb.VoteReply{Id: request.Id}

	if r.currentTerm == request.Term &&
		(r.votedFor == "" || r.votedFor == request.CandidateName) &&
		(request.LastLogTerm > lastLogTerm ||
			(request.LastLogTerm == lastLogTerm && request.LastLogIndex >= lastLogIndex)) {
		r.logger.Debugf("[%s] r.currentTerm= %v; request.Term= %v; r.votedFor= %v; request.CandidateName= %v;"+
			" request.LastLogTerm= %v; lastLogTerm=%v; request.LastLogIndex=%v; lastLogIndex=%v;",
			r.Server().Name, r.currentTerm, request.Term, r.votedFor, request.CandidateName, request.LastLogTerm,
			lastLogTerm, request.LastLogIndex, lastLogIndex)
		reply.VoteGranted = true
		r.votedFor = request.CandidateName
		r.electionResetEvent = time.Now()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = r.currentTerm
	reply.CandidateName = r.server.Name

	err := r.persistToStorage(nil)
	if err != nil {
		return nil, err
	}

	r.logger.Debugf("[%s] [%s] Replying RequestVote Term = %v; VoteGranted = %v; CandidateName = %s",
		r.Server().Name, reply.Id, reply.Term, reply.VoteGranted, reply.CandidateName)

	return reply, nil
}

func (r *raft) startLeader() {
	r.state = Leader

	if r.IsTestMode() {
		if r.ev.LeaderChangeChan != nil {
			r.ev.LeaderChangeChan <- r.Server()
		}
	}

	for _, peer := range r.Peers() {
		r.nextIndex[peer.Name] = int64(len(r.Log()))
		r.matchIndex[peer.Name] = -1
	}
	r.logger.Debugf("[%s] becomes Leader; term=%d, nextIndex=%v, matchIndex=%v;", r.Server().Name, r.currentTerm,
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
				r.logger.Debugf("[%s] Stopping startLeader", r.server.Name)
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

func (r *raft) processAEReply(resp interface{}, savedCurrentTerm int64, peer *pb.Peer, ni int64, entries []*pb.WalEntry) {
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
			r.logger.Debugf("[%s] [%s] AppendEntries reply from %s success: nextIndex := %v, "+
				"matchIndex := %v; commitIndex := %d", r.Server().Name, reply.Id, peer.Name, r.nextIndex,
				r.matchIndex, r.commitIndex)
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
			r.logger.Debugf("[%s] [%s] AppendEntries reply from %s !success: nextIndex := %d",
				r.Server().Name, reply.Id, peer.Name, ni-1)
		}
	}
}

func (r *raft) createAERequest(peer *pb.Peer, savedCurrentTerm int64, id string) (*pb.AppendEntryRequest, int64, []*pb.WalEntry) {
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
		Leader:       r.server,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: r.commitIndex,
		NextIndex:    nextIndex,
		Id:           id,
	}
	return &request, nextIndex, entries
}

func (r *raft) leaderSendAEs() {
	r.mutex.Lock()
	savedCurrentTerm := r.currentTerm
	r.mutex.Unlock()

	id := uuid.New().String()

	for _, peer := range r.Peers() {
		r.wg.Add(1)
		go func(peer *pb.Peer) {
			defer r.wg.Done()

			request, nextIndex, entries := r.createAERequest(peer, savedCurrentTerm, id)

			resp, err := r.rpcServer.Send(peer, "RPC.AppendEntries", request)
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

	r.logger.Debugf("[%s][%s][%s] Received AppendEntries from [%s], current Term = %v", r.Server().Name, r.state,
		request.Id, request.Leader.Name, r.currentTerm)

	if request.Term > r.currentTerm {
		r.logger.Debugf("[%s][%s][%s] request.Term = %v; currentTerm = %v; term out of date in AppendEntries",
			r.Server().Name, r.state, request.Id, request.Term, r.currentTerm)
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

				err := r.persistToStorage(request.Entries[newEntriesIndex:])
				if err != nil {
					return nil, err
				}
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
	r.logger.Debugf("[%s] [%s] [%s] AppendEntries reply: Term = %v; ConflictTerm =%v; ConflictIndex =%v; "+
		"Success =%v;", r.Server().Name, r.state, request.Id, reply.Term, reply.ConflictTerm,
		reply.ConflictIndex, reply.Success)
	return reply, nil
}

//ErrorNotALeader error returned where no server leader
var ErrorNotALeader = errors.New("now a leader")

func (r *raft) submit(command interface{}, commandID string, cmdType pb.CommandType) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != Leader {
		panic(ErrorNotALeader)
		return ErrorNotALeader
	}

	pv, ok := command.(proto.Message)
	if !ok {
		panic(fmt.Sprintf("%v is not proto.Message", pv))
	}

	any, err := anypb.New(pv)
	if err != nil {
		return err
	}

	entry := &pb.WalEntry{
		Command: any,
		Term:    r.currentTerm,
		Id:      commandID,
		CmdType: cmdType,
	}
	r.log = append(r.log, entry)

	err = r.persistToStorage([]*pb.WalEntry{entry})
	if err != nil {
		return err
	}

	r.triggerAEChan <- struct{}{}
	return nil
}

func (r *raft) SubmitAndWait(request interface{}, ID string, cmdType pb.CommandType) error {
	err := r.AddCommandListener(ID)
	if err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		defer r.wg.Done()
		errChan <- r.WaitForCommandCompletion(ID)
	}()

	r.logger.Debugf("XXXXX [%s][%s] Request Id= [%s]", r.Server().Name, r.State(), ID)

	err = r.submit(request, ID, cmdType)
	if err != nil {
		return err
	}
	r.wg.Add(1)
	err = <-errChan
	return err
}

func (r *raft) becomeFollower(term int64) {
	r.logger.Debugf("[%s] becomes Follower with term=%d; log=%v", r.Server().Name, term, r.log)
	r.state = Follower
	r.currentTerm = term
	r.votedFor = ""
	r.electionResetEvent = time.Now()

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runElectionTimer()
		r.logger.Debug("runElectionTimer::becomeFollower")
	}()
}

func (r *raft) applyCommittedEntry(entry *pb.WalEntry) {
	err := r.StorageManager().Apply(entry)
	if c, ok := r.commitsChanMap[entry.Id]; ok {
		c <- err
	}
}

func (r *raft) commitChanSender() {
	for range r.newCommitReadyChan {
		// Find which entries we have to Apply.
		r.mutex.Lock()
		savedLastApplied := r.lastApplied
		var entries []*pb.WalEntry
		if r.commitIndex > r.lastApplied {
			entries = r.log[r.lastApplied+1 : r.commitIndex+1]
			r.lastApplied = r.commitIndex
		}
		r.mutex.Unlock()

		r.logger.Debugf("commitChanSender entries=%v, savedLastApplied=%d", entries, savedLastApplied)

		for _, entry := range entries {
			r.applyCommittedEntry(entry)
		}
	}
	r.logger.Debug("commitChanSender done")
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

func (r *raft) persistToStorage(entries []*pb.WalEntry) error {
	err := r.wal.WriteState(r.currentTerm, r.votedFor)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err := r.wal.Append(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func intMin(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
