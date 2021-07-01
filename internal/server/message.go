package server

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

type Request struct {
	Id        OperationId `json:"id"`
	Arguments interface{} `json:"arguments"`
}

type Response struct {
	Id   string      `json:"id"`
	Resp interface{} `json:"response"`
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

func SendMessage(srv *MemberServer, request *Request) (*Response, error) {
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("StatusCode is not OK: %v. Body: %v ", resp.StatusCode, string(body))
	}
	if err != nil {
		return nil, err
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
