package server

import (
	"context"
	"errors"

	pb "github.com/davinash/yados/internal/proto/gen"
)

var (
	//ErrorPeerAlreadyExists error if the peer with same name already exists
	ErrorPeerAlreadyExists = errors.New("peer with this name already exists in cluster")
)

func (srv *server) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	return srv.Raft().RequestVotes(ctx, request)
}

func (srv *server) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	return srv.Raft().AppendEntries(ctx, request)
}

func (srv *server) AddMember(ctx context.Context, newPeer *pb.NewPeerRequest) (*pb.NewPeerReply, error) {
	EmptyNewMemberReply := &pb.NewPeerReply{Id: newPeer.Id}
	srv.logger.Debugf("[%s] Received AddMember", newPeer.Id)
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	// 1. check if the member with this name already exists
	for _, peer := range srv.Peers() {
		if peer.Name == newPeer.NewPeer.Name {
			return EmptyNewMemberReply, ErrorPeerAlreadyExists
		}
	}
	// Add new member
	err := srv.Raft().AddPeer(&pb.Peer{
		Name:    newPeer.NewPeer.Name,
		Address: newPeer.NewPeer.Address,
		Port:    newPeer.NewPeer.Port,
	})
	if err != nil {
		return nil, err
	}

	return EmptyNewMemberReply, nil
}
