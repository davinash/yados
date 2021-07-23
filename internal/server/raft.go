package server

import (
	"context"
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
)

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
	Peers() map[string]*pb.Member
	SetPeers(peers map[string]*pb.Member)
	AddPeer(peer *pb.Member)
	Term() uint64

	QuorumSize() int
	MemberCount() int

	Leader() string

	ProcessRequestVote(*pb.VoteRequest) (*pb.VoteRequestReply, error)
}

type raft struct {
	logger          *logrus.Entry
	mutex           sync.RWMutex
	state           string
	stopped         chan bool
	wg              sync.WaitGroup
	electionTimeout time.Duration
	name            string
	peers           map[string]*pb.Member
	isTestMode      bool
	leader          string
	currentTerm     uint64
	votedFor        string
	serverName      string
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
		peers:           make(map[string]*pb.Member),
		serverName:      args.ServerName,
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

func (r *raft) Name() string {
	return r.name
}

func (r *raft) SetName(name string) {
	r.name = name
}

func (r *raft) Peers() map[string]*pb.Member {
	return r.peers
}

func (r *raft) SetPeers(peers map[string]*pb.Member) {
	r.peers = peers
}

func (r *raft) AddPeer(peer *pb.Member) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.peers[peer.Name] = peer
}

func (r *raft) Leader() string {
	return r.leader
}

func (r *raft) Term() uint64 {
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

func (r *raft) QuorumSize() int {
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
				go func(peer *pb.Member) {
					conn, p, err := GetPeerConn(peer.Address, peer.Port)
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
		if votesGranted == r.QuorumSize() {
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
		case <-timeoutChan:
			doVote = true
		}
	}
}

func (r *raft) executeLeaderLoop() {
	//for r.State() == Leader {
	//	select {
	//	case <-r.stopped:
	//		r.SetState(Stopped)
	//		return
	//	}
	//}
	r.logger.Debugf("Returning executeLeaderLoop :: stop")
}

func (r *raft) promotable() bool {
	return true
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
