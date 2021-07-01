package server

import "fmt"

func Join(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Adding new member ")

	newPeer, ok := args.(MemberServer)
	if !ok {
		return nil, fmt.Errorf("internal error, invalid type received by join operation")
	}

	// 1. Check if the peer already exists
	// 2. If yes see if the cluster Name is the same
	if peer, ok := server.peers[newPeer.Name]; ok {
		return nil, fmt.Errorf("server with Name %s already exists in cluster %s", peer.Name)
	}

	server.peers[newPeer.Name] = &MemberServer{
		Port:    newPeer.Port,
		Address: newPeer.Address,
		Name:    newPeer.Name,
	}
	return nil, nil
}
