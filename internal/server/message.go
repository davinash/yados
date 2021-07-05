package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type OperationId int

const (
	AddNewMember OperationId = iota
	PutObject
	GetObject
	DeleteObject
	CreateStoreInCluster
	DeleteStoreFromCluster
	StopServer
	AddNewMemberEx
)

type Request struct {
	Id        OperationId `json:"id"`
	Arguments interface{} `json:"arguments"`
}

type Response struct {
	Id    string      `json:"id"`
	Resp  interface{} `json:"response"`
	Error string      `json:"error"`
}

type JoinMember struct {
	Port        int    `json:"Port"`
	Address     string `json:"Address"`
	ClusterName string `json:"ClusterName"`
	Name        string `json:"Name"`
}

type JoinMemberResp struct {
	Port        int    `json:"Port"`
	Address     string `json:"Address"`
	ClusterName string `json:"ClusterName"`
	Name        string `json:"Name"`
}

type StopMember struct {
}

type StopMemberResp struct {
}

func SendMessage(srv *MemberServer, request *Request, logger *logrus.Entry) (*Response, error) {
	url := fmt.Sprintf("http://%s:%d/message", srv.Address, srv.Port)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Errorf("Unmarshling error -> %v", err)
		return nil, err
	}
	if result.Error != "" {
		return nil, fmt.Errorf(result.Error)
	}
	return &result, nil
}

func BroadcastMessage(server *Server, request *Request, logger *logrus.Entry) ([]*Response, error) {
	logger.Debug("Sending Broadcast message")
	allResponses := make([]*Response, 0)
	// Send message to all the members
	for _, srv := range server.peers {
		logger.Debugf("Sending to %s:%d", srv.Address, srv.Port)
		resp, err := SendMessage(srv, request, logger)
		if err != nil {
			return nil, err
		}
		allResponses = append(allResponses, resp)
	}
	return allResponses, nil
}
