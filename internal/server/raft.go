package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	pb "github.com/davinash/yados/internal/proto/gen"

	"github.com/sirupsen/logrus"
)

const (
	//Stopped stopped state
	Stopped = "stopped"
	//Initialized Initialized state
	Initialized = "initialized"
	//Follower Follower state
	Follower = "follower"
	//Candidate Candidate state
	Candidate = "candidate"
	//Leader Leader state
	Leader = "leader"
)

const (
	//DefaultElectionTimeout default election timeout
	DefaultElectionTimeout = 150 * time.Millisecond
	//DefaultHeartbeatInterval default heartbeat interval
	DefaultHeartbeatInterval = 50 * time.Millisecond
)

// ErrorNotLeader = errors.New("this is not a leader") error for returning if not leader
var ErrorNotLeader = errors.New("this is not a leader")

//Raft raft interface
type Raft interface {
	Start() error
	Logger() *logrus.Entry
	SetLogger(logger *logrus.Entry)
	Running() bool
	Initialize() error
	SetState(state string)
	State() string
	Stop() error
	ElectionTimeout() time.Duration
	SetElectionTimeout(duration time.Duration)
	Name() string
	SetName(name string)
	Peers() map[string]*Peer
	SetPeers(peers map[string]*Peer)
	AddPeer(peer *Peer)
	Term() uint64

	IsMajority() int
	MemberCount() int

	Leader() string
	CommandChan() chan *Command

	ProcessRequestVote(*pb.VoteRequest) (*pb.VoteRequestReply, error)

	Log() Log
	SetLog(log Log)
}

//Command represents the internal command for either follower, candidate or leader loop
type Command struct {
	Args      interface{}
	Response  interface{}
	errorChan chan error
}

//Peer add later
type Peer struct {
	stopChan          chan bool
	member            *pb.Member
	heartbeatInterval time.Duration
}

type raft struct {
	logger          *logrus.Entry
	mutex           sync.RWMutex
	state           string
	stopped         chan bool
	wg              sync.WaitGroup
	electionTimeout time.Duration
	name            string
	peers           map[string]*Peer
	isTestMode      bool
	leader          string
	currentTerm     uint64
	votedFor        string
	serverName      string
	commandChan     chan *Command
	log             Log
}

//RftArgs Arguments to create new raft instance
type RftArgs struct {
	Logger     *logrus.Entry
	Name       string
	ServerName string
	IsTestMode bool
}

//NewRaftInstance creates new instance of raft
func NewRaftInstance(args *RftArgs) (Raft, error) {
	r := &raft{logger: args.Logger,
		electionTimeout: DefaultElectionTimeout,
		state:           Stopped,
		peers:           make(map[string]*Peer),
		serverName:      args.ServerName,
		commandChan:     make(chan *Command),
		log:             NewLog(),
	}
	if args.IsTestMode {
		r.isTestMode = true
	}
	return r, nil
}

func (r *raft) Logger() *logrus.Entry {
	return r.logger
}

func (r *raft) SetLogger(logger *logrus.Entry) {
	r.logger = logger
}

func (r *raft) Log() Log {
	return r.log
}

func (r *raft) SetLog(log Log) {
	r.log = log
}

func (r *raft) Name() string {
	return r.name
}

func (r *raft) SetName(name string) {
	r.name = name
}

func (r *raft) Peers() map[string]*Peer {
	return r.peers
}

func (r *raft) SetPeers(peers map[string]*Peer) {
	r.peers = peers
}

func (r *raft) CommandChan() chan *Command {
	return r.commandChan
}

func (r *raft) AddPeer(peer *Peer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.peers[peer.member.Name] = peer
}

func (r *raft) Leader() string {
	return r.leader
}

func (r *raft) Term() uint64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.currentTerm
}

func (r *raft) Running() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.state != Stopped && r.state != Initialized
}

func (r *raft) Initialize() error {
	r.state = Initialized
	return nil
}

func (r *raft) SetState(state string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.state = state

	// Increment current term, vote for self.
	r.currentTerm++
	r.votedFor = r.name

	// Update state and leader.
	r.state = state
	r.logger.Debugf("Current State %s", r.state)
	if state == Leader {
		r.leader = r.Name()
	}
}

func (r *raft) State() string {
	return r.state
}

func (r *raft) ElectionTimeout() time.Duration {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.electionTimeout
}

func (r *raft) SetElectionTimeout(duration time.Duration) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.electionTimeout = duration
}

func (r *raft) MemberCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.peers) + 1
}

func (r *raft) IsMajority() int {
	return (r.MemberCount() / 2) + 1
}

func (r *raft) Start() error {
	if r.Running() {
		return fmt.Errorf("raft instance already running[%v]", r.state)
	}
	err := r.Initialize()
	if err != nil {
		return err
	}

	r.stopped = make(chan bool)
	r.SetState(Follower)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.startRaftLoop()
	}()
	r.logger.Debugf("Started raft ...")
	return nil
}

func (r *raft) startRaftLoop() {
	//r.logger.Debugf("Starting Loop for Raft")
	state := r.State()

	for state != Stopped {
		switch state {
		case Follower:
			r.executeFollowerLoop()
		case Candidate:
			r.executeCandidateLoop()
		case Leader:
			r.executeLeaderLoop()
		}
		state = r.State()
	}
}

func (r *raft) executeFollowerLoop() {
	//since := time.Now()
	//electionTimeout := r.ElectionTimeout()
	timeoutChan := getRandomTimeout(r.ElectionTimeout(), r.ElectionTimeout()*2)

	for r.State() == Follower {

		update := false
		select {
		case <-r.stopped:
			r.SetState(Stopped)
			return
		case <-timeoutChan:
			if r.promotable() {
				r.SetState(Candidate)
			} else {
				update = true
			}
		case c := <-r.commandChan:
			var err error
			switch c.Args.(type) {
			case *pb.VoteRequest:
				c.Response, err = r.ProcessRequestVote(c.Args.(*pb.VoteRequest))
			case *pb.AppendEntriesRequest:
			default:
				err = ErrorNotLeader
			}
			c.errorChan <- err
		}
		if update {
			//since = time.Now()
			timeoutChan = getRandomTimeout(r.ElectionTimeout(), r.ElectionTimeout()*2)
		}
	}
}

func (r *raft) executeCandidateLoop() {
	r.leader = ""
	doVote := true
	votesGranted := 0
	var timeoutChan <-chan time.Time
	var respChan chan *pb.VoteRequestReply

	for r.State() == Candidate {
		if doVote {
			// Increment current term, vote for self.
			r.currentTerm++
			r.votedFor = r.name

			// Request for a vote to peers
			respChan = make(chan *pb.VoteRequestReply, len(r.peers))
			for _, peer := range r.peers {
				go func(peer *Peer) {
					conn, p, err := GetPeerConn(peer.member.Address, peer.member.Port)
					if err != nil {
						return
					}
					voteRequestReply, err := p.RequestForVote(context.Background(), &pb.VoteRequest{
						Term:          r.currentTerm,
						StoreName:     r.name,
						CandidateName: r.serverName,
					})
					if err != nil {
						r.logger.Warnf("RequestForVote Failed, Error = %v", err)
						return
					}

					err = conn.Close()
					if err != nil {
						r.logger.Warnf("failed to close connection, error = %v", err)
					}
					respChan <- voteRequestReply
				}(peer)
			}

			votesGranted = 1
			timeoutChan = getRandomTimeout(r.ElectionTimeout(), r.ElectionTimeout()*2)
			doVote = false

		} // doVote loop finishes here

		// If we received enough votes then stop waiting for more votes.
		// And return from the candidate loop
		if votesGranted == r.IsMajority() {
			r.SetState(Leader)
			return
		}

		select {
		case <-r.stopped:
			r.logger.Debugf("Received executeCandidateLoop :: stop")
			r.SetState(Stopped)
			return
		case resp := <-respChan:
			if success := r.processVoteResponse(resp); success {
				r.logger.Debugf("votesGranted")
				votesGranted++
			}
		case c := <-r.commandChan:
			var err error
			switch c.Args.(type) {
			case *pb.VoteRequest:
				c.Response, err = r.ProcessRequestVote(c.Args.(*pb.VoteRequest))
			case *pb.AppendEntriesRequest:
			}
			c.errorChan <- err
		case <-timeoutChan:
			doVote = true
		}
	}
}

func (r *raft) executeLeaderLoop() {
	for _, peer := range r.peers {
		r.startHeartbeat(peer)
	}

	for r.State() == Leader {
		select {
		case <-r.stopped:
			r.SetState(Stopped)
			return
		case c := <-r.commandChan:
			var err error
			switch c.Args.(type) {
			case *pb.VoteRequest:
				c.Response, err = r.ProcessRequestVote(c.Args.(*pb.VoteRequest))
			case *pb.AppendEntriesRequest:
			default:
				err = ErrorNotLeader
			}
			c.errorChan <- err
		}
	}
	r.logger.Debugf("Returning executeLeaderLoop :: stop")
}

func (r *raft) promotable() bool {
	return r.Log().CurrentIndex() > 0
}

func getRandomTimeout(min time.Duration, max time.Duration) <-chan time.Time {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d, delta := min, max-min
	if delta > 0 {
		d += time.Duration(r.Int63n(int64(delta)))
	}
	return time.After(d)
}

func (r *raft) Stop() error {
	if r.State() == Stopped {
		return nil
	}
	close(r.stopped)
	r.wg.Wait()
	r.SetState(Stopped)
	r.logger.Debugf("Stopped Raft Instance")
	return nil
}

//HeartBeat add later
func (p *Peer) HeartBeat(c chan bool) {
	stopChan := p.stopChan
	c <- true
	ticker := time.Tick(p.heartbeatInterval)
	for {
		select {
		case <-stopChan:
			return
		case <-ticker:
		}
	}
}

func (r *raft) startHeartbeat(peer *Peer) {
	peer.stopChan = make(chan bool)
	c := make(chan bool)
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		peer.HeartBeat(c)
	}()
	<-c
}
