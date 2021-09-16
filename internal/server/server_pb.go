package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davinash/yados/internal/store"
	"google.golang.org/grpc"

	"github.com/davinash/yados/internal/raft"

	"github.com/sirupsen/logrus"

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

func (srv *server) RemovePeer(ctx context.Context, request *pb.RemovePeerRequest) (*pb.RemovePeerReply, error) {
	EmptyRemovePeerReply := &pb.RemovePeerReply{Id: request.Id}
	srv.logger.Debugf("[%s] [%s] Received RemovePeer", srv.Name(), request.Id)

	err := srv.Raft().RemovePeer(request)
	if err != nil {
		return EmptyRemovePeerReply, err
	}
	return EmptyRemovePeerReply, nil
}

func (srv *server) PeerStatus(ctx context.Context, request *pb.StatusRequest) (*pb.StatusReply, error) {
	reply := &pb.StatusReply{
		Id:       request.Id,
		IsLeader: false,
	}

	reply.Server = srv.Self()
	reply.Status = srv.Raft().State().String()
	if srv.Raft().State() == raft.Leader {
		reply.IsLeader = true
	}
	return reply, nil
}

func (srv *server) ClusterStatus(ctx context.Context, request *pb.ClusterStatusRequest) (*pb.ClusterStatusReply, error) {
	reply := &pb.ClusterStatusReply{Id: request.Id}
	for _, peer := range srv.Peers() {
		args := pb.StatusRequest{}
		status, err := srv.RPCServer().Send(peer, "RPC.PeerStatus", &args)
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

//ErrorStoreAlreadyExists error if store with this name already exists
var ErrorStoreAlreadyExists = errors.New("store with this name already exists")

func (srv *server) CreateStore(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	srv.logger.Debugf("[%s] [%s] XXXXX Received request CreateStore id=%v Name=%v, Type=%v",
		srv.Name(), srv.State(), request.Id, request.Name, request.Type)
	reply := &pb.StoreCreateReply{}
	// Check if store with name already exists
	if _, ok := srv.StoreManager().Stores()[request.Name]; ok {
		return reply, ErrorStoreAlreadyExists
	}

	err := srv.Raft().SubmitAndWait(request, request.Id, pb.CommandType_CreateStore)
	if err != nil {
		return reply, err
	}
	reply.Msg = fmt.Sprintf("Store = %s, Type = %s created", request.Name, request.Type.String())
	return reply, nil
}

//ErrorStoreDoesExists error if store does not exists
var ErrorStoreDoesExists = errors.New("store does not exists")

func (srv *server) DeleteStore(ctx context.Context, request *pb.StoreDeleteRequest) (*pb.StoreDeleteReply, error) {
	reply := &pb.StoreDeleteReply{}
	if _, ok := srv.StoreManager().Stores()[request.StoreName]; !ok {
		reply.Error = ErrorStoreDoesExists.Error()
		return reply, ErrorStoreDoesExists
	}
	if srv.logger.IsLevelEnabled(logrus.DebugLevel) {
		marshal, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		srv.logger.Debugf("[%s] Submitting request to raft engine %s", request.Id, string(marshal))
	}
	err := srv.Raft().SubmitAndWait(request, request.Id, pb.CommandType_DeleteStore)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

func (srv *server) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutReply, error) {
	reply := &pb.PutReply{}
	// Check if store with name exists
	if _, ok := srv.StoreManager().Stores()[request.StoreName]; !ok {
		reply.Error = ErrorStoreDoesExists.Error()
		return reply, ErrorStoreDoesExists
	}

	if srv.logger.IsLevelEnabled(logrus.DebugLevel) {
		marshal, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		srv.logger.Debugf("[%s] Submitting request to raft engine %s", request.Id, string(marshal))
	}

	err := srv.Raft().SubmitAndWait(request, request.Id, pb.CommandType_Put)
	if err != nil {
		return reply, err
	}

	return reply, nil
}

func (srv *server) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetReply, error) {
	reply := &pb.GetReply{}
	// Check if store with name exists
	if _, ok := srv.StoreManager().Stores()[request.StoreName]; !ok {
		reply.Value = ErrorStoreDoesExists.Error()
		return reply, ErrorStoreDoesExists
	}

	value := (srv.StoreManager().Stores()[request.StoreName].(store.KVStore)).Get(request)
	reply.Value = value

	return reply, nil
}

func (srv *server) ListStores(ctx context.Context, request *pb.ListStoreRequest) (*pb.ListStoreReply, error) {
	reply := &pb.ListStoreReply{}
	for k := range srv.StoreManager().Stores() {
		reply.Name = append(reply.Name, k)
	}
	return reply, nil
}

func (srv *server) ExecuteQuery(ctx context.Context, request *pb.ExecuteQueryRequest) (*pb.ExecuteQueryReply, error) {
	reply := &pb.ExecuteQueryReply{Id: request.Id}

	// Check if store with name exists
	if _, ok := srv.StoreManager().Stores()[request.StoreName]; !ok {
		return reply, ErrorStoreDoesExists
	}

	err := srv.Raft().SubmitAndWait(request, request.Id, pb.CommandType_SqlDDL)
	return reply, err
}

func (srv *server) Query(ctx context.Context, request *pb.QueryRequest) (*pb.QueryReply, error) {
	reply := &pb.QueryReply{}
	if _, ok := srv.StoreManager().Stores()[request.StoreName]; !ok {
		return reply, ErrorStoreDoesExists
	}

	resp, err := (srv.StoreManager().Stores()[request.StoreName].(store.SQLStore)).Query(request)
	return resp, err
}

func (srv *server) AddPeers(ctx context.Context, peers *pb.AddPeersRequest) (*pb.AddPeersReply, error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	for _, p := range peers.GetPeers() {
		if p.Name != srv.Name() {
			srv.Raft().AddPeer(p)
		}
	}
	return &pb.AddPeersReply{}, nil
}

func (srv *server) GetLeader(ctx context.Context, request *pb.GetLeaderRequest) (*pb.GetLeaderReply, error) {
	controllers := srv.Controller()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", controllers[0].address, controllers[0].port), grpc.WithInsecure())
	if err != nil {
		return &pb.GetLeaderReply{}, fmt.Errorf("failed to connect with controller[%s:%d], "+
			"error = %w", controllers[0].address, controllers[0].port, err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			srv.logger.Warnf("failed to close the connection with controller, Error = %v", err)
		}
	}(conn)
	controller := pb.NewControllerServiceClient(conn)
	leader, err := controller.GetLeader(ctx, request)
	if err != nil {
		return &pb.GetLeaderReply{}, err
	}
	return leader, nil
}
