package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type MemberServer struct {
	Port        int
	Address     string
	clusterName string
	name        string
}

type Server struct {
	self       *MemberServer
	listener   net.Listener
	OSSignalCh chan os.Signal
	peers      map[string]*MemberServer
	isTestMode bool
}

func CreateNewServer(name string, address string, port int, clusterName string,
	withPeer bool, peerAddress string, peerPort int) (*Server, error) {
	server := Server{
		self: &MemberServer{
			Port:        port,
			Address:     address,
			clusterName: clusterName,
			name:        name,
		},
		peers:      map[string]*MemberServer{},
		isTestMode: false,
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

func (server *Server) EnableTestMode() {
	server.isTestMode = true
}

func (server *Server) startHttpServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.self.Address, server.self.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: Error = %v", err)
	}
	server.listener = listener
	go func() {
		http.Serve(server.listener, SetupRouter(server))
	}()
	return nil
}

func (server *Server) BroadcastMessage(request *Request) ([]*Response, error) {
	allResponses := make([]*Response, 0)
	// Send message to all the members
	for _, srv := range server.peers {
		resp, err := SendMessage(srv, request)
		if err != nil {
			return nil, err
		}
		allResponses = append(allResponses, resp)
	}
	return allResponses, nil
}

func (server *Server) PostInit() error {
	log.Println("Performing Post Initialization ...")
	_, err := server.BroadcastMessage(&Request{
		Id: AddNewMember,
		Arguments: MemberServer{
			Port:        server.self.Port,
			Address:     server.self.Address,
			clusterName: server.self.clusterName,
			name:        server.self.name,
		},
	})
	if err != nil {
		return err
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
	log.Println("Stopping the server ...")
	if !server.isTestMode {
		if err := server.listener.Close(); err != nil {
			log.Printf("failed to stop the http server, Error = %v\n", err)
			return err
		}
		os.Exit(0)
	}
	return nil
}

func (server *Server) Status() error {
	return nil
}
