package server

//ListMemberFn list all the members in the cluster
func ListMemberFn(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Performing ListMemberFn")
	members := make([]MemberServer, 0)
	// Add all the neighbor
	for _, srv := range server.peers {
		members = append(members, MemberServer{
			Port:    srv.Port,
			Address: srv.Address,
			Name:    srv.Name,
		})
	}
	// Add self
	members = append(members, MemberServer{
		Port:    server.self.Port,
		Address: server.self.Address,
		Name:    server.self.Name,
	})

	resp := Response{
		ID:    "",
		Resp:  members,
		Error: "",
	}
	return &resp, nil
}
