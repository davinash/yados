package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type MemberServer struct {
	port        int
	address     string
	clusterName string
}

type Server struct {
	name        string
	port        int
	listener    net.Listener
	OSSignalCh  chan os.Signal
	address     string
	clusterName string
	members     map[string]*MemberServer
}

func GetExistingServer() (*Server, error) {
	server := Server{}
	return &server, nil
}

func CreateNewServer(name string, address string, port int, clusterName string) (*Server, error) {
	server := Server{
		name:        name,
		address:     address,
		port:        port,
		clusterName: clusterName,
		members:     map[string]*MemberServer{},
	}
	server.OSSignalCh = make(chan os.Signal, 1)
	return &server, nil
}

func (server *Server) HandleSignal() {
	for {
		select {
		case _ = <-server.OSSignalCh:
			log.Println("Exiting ... ")
			server.Stop()
			os.Exit(0)
		}
	}
}

func (server *Server) startHttpServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
	if err != nil {
		return fmt.Errorf("failed to start server: Error = %v", err)
	}
	server.listener = listener
	go func() {
		http.Serve(server.listener, setupRouter(server))
	}()
	return nil
}

func (server *Server) SendMessage(srv *MemberServer, request *Request) (*Response, error) {
	url := fmt.Sprintf("%s:%d/message", srv.address, srv.port)
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
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
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (server *Server) BroadcastMessage(request *Request) ([]*Response, error) {
	allResponses := make([]*Response, 0)
	for _, srv := range server.members {
		resp, err := server.SendMessage(srv, request)
		if err != nil {
			return nil, err
		}
		allResponses = append(allResponses, resp)
	}
	return allResponses, nil
}

func (server *Server) PostInit() error {
	log.Println("Performing Post Initialization ...")
	if member, ok := server.members[server.name]; ok {
		if member.clusterName == server.clusterName {
			return fmt.Errorf("server with name %s already exists in cluster %s", server.name, server.clusterName)
		}
	}
	server.members[server.name] = &MemberServer{
		port:        server.port,
		address:     server.address,
		clusterName: server.clusterName,
	}
	return nil
}

func (server *Server) StartAndWait() error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.startHttpServer()
	if err != nil {
		return err
	}
	err = server.PostInit()
	if err != nil {
		server.Stop()
		return err
	}
	<-server.OSSignalCh
	return nil
}

func (server *Server) Stop() error {
	err := server.listener.Close()
	if err != nil {
		log.Printf("failed to stop the http server, Error = %v\n", err)
		return err
	}
	return nil
}

func (server *Server) Status() error {
	return nil
}
