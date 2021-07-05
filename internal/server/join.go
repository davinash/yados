package server

import (
	"encoding/json"
	"fmt"
)

//Join implements message for joining a new member in the cluster
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
		ID:        AddNewMemberEx,
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

//JoinEx Adds a new member in the local dictionary
func JoinEx(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Adding New Member from existing member")

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
