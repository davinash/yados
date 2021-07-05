package server

func ListAllMembers(args interface{}, server *Server) (*Response, error) {
	server.logger.Info("Performing ListAllMembers")
	members := make([]MemberServer, 0)
	for _, srv := range server.peers {
		members = append(members, MemberServer{
			Port:    srv.Port,
			Address: srv.Address,
			Name:    srv.Name,
		})
	}
	members = append(members, MemberServer{
		Port:    server.self.Port,
		Address: server.self.Address,
		Name:    server.self.Name,
	})

	resp := Response{
		Id:    "",
		Resp:  members,
		Error: "",
	}
	server.logger.Info("-----> %+v", resp)
	return &resp, nil
}
