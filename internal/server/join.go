package server

import "log"

func Join(args interface{}) Response {
	log.Println("Adding new member ")
	return Response{
		Id: "",
		Resp: JoinMemberResp{
			Port:        0,
			Address:     "",
			ClusterName: "",
			Name:        "",
		},
		Err: nil,
	}

	//if member, ok := server.members[server.name]; ok {
	//	if member.clusterName == server.clusterName {
	//		return fmt.Errorf("server with name %s already exists in cluster %s", server.name, server.clusterName)
	//	}
	//}
	//server.members[server.name] = &MemberServer{
	//	port:        server.port,
	//	address:     server.address,
	//	clusterName: server.clusterName,
	//}
}
