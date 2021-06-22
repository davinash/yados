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

func (server *Server) Start() error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if member, ok := server.members[server.name]; ok {
		if member.clusterName == server.clusterName {
			return fmt.Errorf("server with name %s already exists in cluster %s", server.name, server.clusterName)
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	server.listener = listener

	go func() {
		err := http.Serve(server.listener, setupRouter(server))
		if err != nil {
			log.Fatalf("Failed to start a server, Error = %v", err)
		}
	}()
	// Make your own entry into members map
	server.members[server.name] = &MemberServer{
		port:        server.port,
		address:     server.address,
		clusterName: server.clusterName,
	}

	for {
		select {
		case _ = <-server.OSSignalCh:
			log.Println("Exiting ... ")
			os.Exit(0)
		}
	}
	return nil
}

func (server *Server) Stop() error {
	return nil
}

func (server *Server) Status() error {
	return nil
}
