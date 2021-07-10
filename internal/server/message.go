package server

//
////OperationID represents the operation id or request
//type OperationID int
//
//const (
//	//AddNewMemberInCluster request type to add new member.proto
//	AddNewMemberInCluster OperationID = iota
//	//PutObject request type to put new object in a store
//	PutObject
//	//GetObject request type to get object from store
//	GetObject
//	//DeleteObject request type to delete object
//	DeleteObject
//	//CreateStoreInCluster request type to create cluster
//	CreateStoreInCluster
//	//DeleteStoreFromCluster request type to delete store
//	DeleteStoreFromCluster
//	//StopServer request type to stop a server
//	StopServer
//	//AddNewMemberEx request type to add new member.proto
//	AddNewMemberEx
//	//ListMembers request type to list all the members
//	ListMembers
//)
//
////Request represents the HTTP request
//type Request struct {
//	ID        OperationID `json:"id"`
//	Arguments interface{} `json:"arguments"`
//}
//
////Response represents response for the HTTP Request
//type Response struct {
//	ID    string      `json:"id"`
//	Resp  interface{} `json:"response"`
//	Error string      `json:"error"`
//}
//
////JoinMember represent the request for adding a new member.proto in the cluster
//type JoinMember struct {
//	Port        int    `json:"Port"`
//	Address     string `json:"Address"`
//	ClusterName string `json:"ClusterName"`
//	Name        string `json:"Name"`
//}
//
////JoinMemberResp response for the adding new member.proto
//type JoinMemberResp struct {
//	Port        int    `json:"Port"`
//	Address     string `json:"Address"`
//	ClusterName string `json:"ClusterName"`
//	Name        string `json:"Name"`
//}
//
////StopMember Request for stopping a member.proto
//type StopMember struct {
//}
//
////StopMemberResp response from the stopping a member.proto request
//type StopMemberResp struct {
//}
//
////SendMessage Sends message to particular member.proto in the cluster
//func SendMessage(srv *MemberServer, request *Request, logger *logrus.Entry) (*Response, error) {
//	url := fmt.Sprintf("http://%s:%d/message", srv.Address, srv.Port)
//
//	requestBytes, err := json.Marshal(request)
//	if err != nil {
//		return nil, err
//	}
//	r, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
//	if err != nil {
//		return nil, err
//	}
//	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//	client := &http.Client{}
//	resp, err := client.Do(r)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	var result Response
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		logger.Errorf("Unmarshling error -> %v", err)
//		return nil, err
//	}
//	if result.Error != "" {
//		return nil, fmt.Errorf(result.Error)
//	}
//	return &result, nil
//}
//
//// BroadcastMessage Sends the message to all the members in the cluster
//func BroadcastMessage(server *YadosServer, request *Request, logger *logrus.Entry) ([]*Response, error) {
//	logger.Debug("Sending Broadcast message")
//	allResponses := make([]*Response, 0)
//	// Send message to all the members
//	for _, srv := range server.peers {
//		logger.Debugf("Sending to %s:%d", srv.Address, srv.Port)
//		resp, err := SendMessage(srv, request, logger)
//		if err != nil {
//			return nil, err
//		}
//		allResponses = append(allResponses, resp)
//	}
//	return allResponses, nil
//}
