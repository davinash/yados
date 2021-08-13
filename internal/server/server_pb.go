package server

import (
	"context"
	"errors"
	"log"
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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

func (srv *server) PeerStatus(ctx context.Context, request *pb.StatusRequest) (*pb.StatusReply, error) {
	reply := &pb.StatusReply{Id: request.Id}
	reply.Server = srv.Self()
	reply.Status = srv.State().String()
	return reply, nil
}

func (srv *server) ClusterStatus(ctx context.Context, request *pb.ClusterStatusRequest) (*pb.ClusterStatusReply, error) {
	reply := &pb.ClusterStatusReply{Id: request.Id}
	for _, peer := range srv.Peers() {
		args := pb.StatusRequest{}
		status, err := srv.Send(peer, "RPC.PeerStatus", &args)
		if err != nil {
			srv.logger.Errorf("failed to send PeerStatus to %s, Error = %v", peer.Name, err)
			return reply, err
		}
		reply.PeerStatus = append(reply.PeerStatus, status.(*pb.StatusReply))
	}
	status, err := srv.PeerStatus(ctx, &pb.StatusRequest{})
	if err != nil {
		return nil, err
	}
	reply.PeerStatus = append(reply.PeerStatus, status)
	return reply, nil
}

func (srv *server) SubmitToRaft(requestBytes []byte, ID string, cmdType pb.CommandType) error {
	err := srv.Raft().AddCommandListener(ID)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.Raft().WaitForCommandCompletion(requestBytes, ID, cmdType)
	}()

	err = srv.Raft().Submit(requestBytes, ID, cmdType)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

//ErrorStoreAlreadyExists error if store with this name already exists
var ErrorStoreAlreadyExists = errors.New("store with this name already exists")

func (srv *server) CreateStore(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	reply := &pb.StoreCreateReply{}
	// Check if store with name already exists
	if _, ok := srv.Stores()[request.Name]; ok {
		reply.Error = ErrorStoreAlreadyExists.Error()
		return reply, ErrorStoreAlreadyExists
	}

	if srv.IsLeader() {
		srv.Logger().Debug("yay this is leader")
		requestBytes, err := proto.Marshal(request)
		if err != nil {
			return reply, err
		}
		srv.logger.Debug("submitting request to raft engine")
		err = srv.SubmitToRaft(requestBytes, request.Id, pb.CommandType_CreateStore)
		if err != nil {
			return reply, err
		}
	} else {
		srv.Logger().Debugf("not a leader, hoping request to %s:%d", srv.leader.Address, srv.leader.Port)
		peerConn, rpcClient, err := GetPeerConn(srv.leader.Address, srv.leader.Port)
		if err != nil {
			return reply, nil
		}
		defer func(peerConn *grpc.ClientConn) {
			err := peerConn.Close()
			if err != nil {
				log.Printf("failed to close the connection, error = %v\n", err)
			}
		}(peerConn)
		_, err = rpcClient.CreateStore(context.Background(), request)
		if err != nil {
			return reply, nil
		}
	}
	return reply, nil
}

func (srv *server) RunCommand(ctx context.Context, command *pb.CommandRequest) (*pb.CommandReply, error) {
	reply := &pb.CommandReply{}
	switch command.CmdType {
	case pb.CommandType_CreateStore:
		var request pb.StoreCreateRequest
		err := proto.Unmarshal(command.Args, &request)
		if err != nil {
			return reply, err
		}
		request.Id = command.Id

		_, err = srv.CreateStore(ctx, &request)
		if err != nil {
			return reply, err
		}
	}
	return reply, nil
}
