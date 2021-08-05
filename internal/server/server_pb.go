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

func (srv *server) RemovePeer(ctx context.Context, request *pb.RemovePeerRequest) (*pb.RemovePeerReply, error) {
	EmptyRemovePeerReply := &pb.RemovePeerReply{Id: request.Id}
	srv.logger.Debugf("[%s] Received RemovePeer", request.Id)

	err := srv.Raft().RemovePeer(request)
	if err != nil {
		return EmptyRemovePeerReply, err
	}
	return EmptyRemovePeerReply, nil
}

func (srv *server) ListMembers(ctx context.Context, request *pb.ListMembersRequest) (*pb.ListMembersReply, error) {
	srv.logger.Debugf("[%s] Received ListMembers", request.Id)
	reply := &pb.ListMembersReply{Id: request.Id}

	peers := srv.Peers()
	for _, p := range peers {
		reply.Peers = append(reply.Peers, p)
	}
	self := srv.Self()
	self.State = srv.State()
	reply.Peers = append(reply.Peers, self)

	return reply, nil
}
