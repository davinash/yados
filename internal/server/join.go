package server

import (
	"encoding/json"
	"fmt"
)

func Join(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Adding new member ")

	marshal, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	var newPeer MemberServer
	err = json.Unmarshal(marshal, &newPeer)
	if err != nil {
		return nil, err
	}
	if peer, ok := server.peers[newPeer.Name]; ok {
		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
	}

	_, err = BroadcastMessage(server, &Request{
		Id:        AddNewMemberEx,
		Arguments: newPeer,
	}, server.logger)
	if err != nil {
		return nil, err
	}
	server.peers[newPeer.Name] = &MemberServer{
		Port:    newPeer.Port,
		Address: newPeer.Address,
		Name:    newPeer.Name,
	}
	return nil, nil
}

func JoinEx(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Adding New Member from existing member")
	//newPeer, ok := args.(MemberServer)
	//if !ok {
	//	return nil, fmt.Errorf("internal error, invalid type received by join operation")
	//}

	marshal, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	var newPeer MemberServer
	err = json.Unmarshal(marshal, &newPeer)
	if err != nil {
		return nil, err
	}

	// 1. Check if the peer already exists
	// 2. If yes see if the cluster Name is the same
	if peer, ok := server.peers[newPeer.Name]; ok {
		return nil, fmt.Errorf("server with Name %s already exists in cluster", peer.Name)
	}

	server.peers[newPeer.Name] = &MemberServer{
		Port:    newPeer.Port,
		Address: newPeer.Address,
		Name:    newPeer.Name,
	}
	return nil, nil
}
