package server

import (
	"fmt"
	"log"
)

func Join(args interface{}, server *Server) (*Response, error) {
	log.Println("Adding new member ")

	newPeer, ok := args.(MemberServer)
	if !ok {
		return nil, fmt.Errorf("internal error, invalid type received by join operation")
	}

	if server.self.clusterName == newPeer.clusterName {
		return nil, fmt.Errorf("cluster name did not match")
	}
	// 1. Check if the peer already exists
	// 2. If yes see if the cluster name is the same
	if peer, ok := server.peers[newPeer.name]; ok {
		return nil, fmt.Errorf("server with name %s already exists in cluster %s", peer.name,
			newPeer.clusterName)
	}

	server.peers[newPeer.name] = &MemberServer{
		Port:        newPeer.Port,
		Address:     newPeer.Address,
		clusterName: newPeer.clusterName,
	}
	return nil, nil
}
