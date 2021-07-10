package server

import (
	"context"
	"fmt"
	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

func (server *YadosServer) AddNewMemberInCluster(ctx context.Context, newPeer *pb.NewMemberRequest) (*pb.NewMemberReply, error) {
	server.logger.Info("Adding new member in the cluster")
	if peer, ok := server.peers[newPeer.Member.Name]; ok {
		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
	}

	for _, p := range server.peers {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.Address, p.Port), grpc.WithInsecure())
		if err != nil {
			_ = conn.Close()
			return &pb.NewMemberReply{}, err
		}
		pb.NewYadosServiceClient(conn)

		_ = conn.Close()
	}

	return &pb.NewMemberReply{}, nil
}

//
////AddNewMemberInClusterFn implements message for joining a new member.proto in the cluster
//func AddNewMemberInClusterFn(args interface{}, server *YadosServer) (*Response, error) {
//	server.logger.Info("Adding new member.proto in cluster")
//
//	marshal, err := json.Marshal(args)
//	if err != nil {
//		return nil, err
//	}
//	var newPeer MemberServer
//	err = json.Unmarshal(marshal, &newPeer)
//	if err != nil {
//		return nil, err
//	}
//	// Check if the member.proto with same name already exists
//	if peer, ok := server.peers[newPeer.Name]; ok {
//		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
//	}
//	// Send message to other peers about this new member.proto
//	r, err := BroadcastMessage(server, &Request{
//		ID:        AddNewMemberEx,
//		Arguments: newPeer,
//	}, server.logger)
//	if err != nil {
//		return nil, err
//	}
//	// Add new member.proto in self
//	server.peers[newPeer.Name] = &MemberServer{
//		Port:    newPeer.Port,
//		Address: newPeer.Address,
//		Name:    newPeer.Name,
//	}
//
//	marshal, err = json.Marshal(r)
//	if err != nil {
//		return nil, err
//	}
//	var resp []*Response
//	err = json.Unmarshal(marshal, &resp)
//	if err != nil {
//		return nil, err
//	}
//	peers := make([]*MemberServer, 0)
//	// Add all peers in the response
//
//	//for _, peer := range resp {
//	//	peers = append(peers, peer.Resp)
//	//}
//
//	// Add self
//	peers = append(peers, server.self)
//
//	return &Response{
//		ID:    "",
//		Resp:  peers,
//		Error: "",
//	}, nil
//}
//
////AddNewMemberExFn Adds a new member.proto in the local dictionary
//func AddNewMemberExFn(args interface{}, server *YadosServer) (*Response, error) {
//	server.logger.Info("Adding New Member from existing member.proto")
//
//	marshal, err := json.Marshal(args)
//	if err != nil {
//		return nil, err
//	}
//	var newPeer MemberServer
//	err = json.Unmarshal(marshal, &newPeer)
//	if err != nil {
//		return nil, err
//	}
//
//	// 1. Check if the peer already exists
//	// 2. If yes see if the cluster Name is the same
//	if peer, ok := server.peers[newPeer.Name]; ok {
//		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
//	}
//
//	server.peers[newPeer.Name] = &MemberServer{
//		Port:    newPeer.Port,
//		Address: newPeer.Address,
//		Name:    newPeer.Name,
//	}
//	return &Response{
//		ID:    "",
//		Resp:  server.self,
//		Error: "",
//	}, nil
//}
