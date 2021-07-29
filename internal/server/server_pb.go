package server

import (
	"context"
	"errors"

	pb "github.com/davinash/yados/internal/proto/gen"
)

var (
	//EmptyVoteReply to return from the vote reply operation
	EmptyVoteReply = &pb.VoteReply{}
	//EmptyAppendEntryReply empty reply from append entry operation
	EmptyAppendEntryReply = &pb.AppendEntryReply{}
	//EmptyNewMemberReply empty reply from Add member operation
	EmptyNewMemberReply = &pb.NewPeerReply{}

	//ErrorPeerAlreadyExists error if the peer with same name already exists
	ErrorPeerAlreadyExists = errors.New("peer with this name already exists in cluster")
)

func (srv *server) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	return srv.Raft().RequestVotes(ctx, request)
}

func (srv *server) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	return srv.Raft().AppendEntries(ctx, request)
}

func (srv *server) AddMember(ctx context.Context, newPeer *pb.Peer) (*pb.NewPeerReply, error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	// 1. check if the member with this name already exists
	for _, peer := range srv.Peers() {
		if peer.Name == newPeer.Name {
			return EmptyNewMemberReply, ErrorPeerAlreadyExists
		}
	}
	// Add new member
	srv.Raft().AddPeer(newPeer)

	return EmptyNewMemberReply, nil
}
