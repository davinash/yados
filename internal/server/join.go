package server

import (
	"fmt"
	"log"
)

func AddNewMemberLocally(server *Server) error {
	if member, ok := server.members[server.self.name]; ok {
		if member.clusterName == server.self.clusterName {
			return fmt.Errorf("server with name %s already exists in cluster %s", server.self.name,
				server.self.clusterName)
		}
	}
	server.members[server.self.name] = &MemberServer{
		port:        server.self.port,
		address:     server.self.address,
		clusterName: server.self.clusterName,
	}
	return nil
}
func Join(args interface{}, server *Server) (*Response, error) {
	log.Println("Adding new member ")
	err := AddNewMemberLocally(server)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
