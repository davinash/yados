package server

import (
	"fmt"
	"log"
)

func Join(args interface{}, server *Server) (*Response, error) {
	log.Println("Adding new member ")

	if member, ok := server.members[server.name]; ok {
		if member.clusterName == server.clusterName {
			return nil, fmt.Errorf("server with name %s already exists in cluster %s", server.name, server.clusterName)
		}
	}
	server.members[server.name] = &MemberServer{
		port:        server.port,
		address:     server.address,
		clusterName: server.clusterName,
	}
	return nil, nil
}
