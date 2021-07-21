package server

import (
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
}

//RftArgs Arguments to create new raft instance
type RftArgs struct {
	Logger *logrus.Entry
	Name   string
}

//NewRaftInstance creates new instance of raft
func NewRaftInstance(args *RftArgs) (Raft, error) {
	return &raft{logger: args.Logger,
		electionTimeout: DefaultElectionTimeout,
		state:           Stopped,
	}, nil
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
	r.peers[peer.Name] = peer
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

	r.logger.Debugf("Starting the Raft Loop")
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.startRaftLoop()
	}()
	r.logger.Debugf("Started raft ...")
	return nil
}

func (r *raft) startRaftLoop() {
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
	for r.State() == Follower {
		select {
		case <-r.stopped:
			r.SetState(Stopped)
			return
		}
	}
}

func (r *raft) executeCandidateLoop() {
	for r.State() == Candidate {
		select {
		case <-r.stopped:
			r.SetState(Stopped)
			return
		}
	}
}

func (r *raft) executeLeaderLoop() {
	for r.State() == Leader {
		select {
		case <-r.stopped:
			r.SetState(Stopped)
			return
		}
	}
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
	r.logger.Debugf("Stopping Raft Instance")
	if r.State() == Stopped {
		return nil
	}
	close(r.stopped)
	r.wg.Wait()

	r.SetState(Stopped)
	return nil
}
