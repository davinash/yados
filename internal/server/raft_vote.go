package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/pkg/errors"
)

var (
	emptyVoteRequestReply = &pb.VoteRequestReply{}
	// ErrorDenyStaleTerm error returned when term is stale
	ErrorDenyStaleTerm = errors.New("Vote Denied, Stale Term")
	//ErrorDenyDuplicate error returned when follower has already voted
	ErrorDenyDuplicate = errors.New("Vote Denied, Duplicate")
)

func (r *raft) processVoteResponse(resp *pb.VoteRequestReply) bool {
	if resp.VoteGranted && resp.Term == r.currentTerm {
		return true
	}
	if resp.Term > r.currentTerm {
		r.updateCurrentTerm(resp.Term, "")
	}
	return false
}

func (r *raft) updateCurrentTerm(term uint64, leaderName string) {
	if term < r.currentTerm {
		r.logger.Debug("update is called when term is not larger than currentTerm")
		panic("update is called when term is not larger than currentTerm")
	}
	// TODO : Stop the heartbeat
	// update the term and clear vote for
	if r.state != Follower {
		r.SetState(Follower)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.currentTerm = term
	r.leader = leaderName
	r.votedFor = ""

}

func (r *raft) ProcessRequestVote(request *pb.VoteRequest) (*pb.VoteRequestReply, error) {
	// If the request is coming from an old term then reject it.
	if request.Term < r.Term() {
		return emptyVoteRequestReply, ErrorDenyStaleTerm
	}
	// If the term of the request peer is larger than this node, update the term
	// If the term is equal and we've already voted for a different candidate then
	// don't vote for this candidate.
	if request.Term > r.Term() {
		r.updateCurrentTerm(request.Term, "")
	} else if r.votedFor != "" && r.votedFor != request.CandidateName {
		return emptyVoteRequestReply, ErrorDenyDuplicate
	}

	// TODO : If the candidate's log is not at least as up-to-date as our last log then don't vote.

	r.votedFor = request.CandidateName
	return &pb.VoteRequestReply{
		Term:        r.currentTerm,
		VoteGranted: true,
	}, nil
}

func (server *server) RequestForVote(ctx context.Context, request *pb.VoteRequest) (*pb.VoteRequestReply, error) {
	command := &Command{
		Args:      request,
		errorChan: make(chan error, 1),
	}
	if s, ok := server.Stores()[request.StoreName]; ok {
		s.RaftInstance().CommandChan() <- command
	}
	err := <-command.errorChan
	if err != nil {
		return emptyVoteRequestReply, err
	}
	return command.Response.(*pb.VoteRequestReply), nil
}
