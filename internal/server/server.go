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

func (server *Server) Start() error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	if member, ok := server.members[server.name]; ok {
		if member.clusterName == server.clusterName {
			return fmt.Errorf("server with name %s already exists in cluster %s", server.name, server.clusterName)
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.address, server.port))
	if err != nil {
		return fmt.Errorf("failed to start server: Error = %v", err)
	}
	server.listener = listener
	go func() {
		http.Serve(server.listener, setupRouter(server))
	}()

	server.members[server.name] = &MemberServer{
		port:        server.port,
		address:     server.address,
		clusterName: server.clusterName,
	}
	return nil
}

func (server *Server) StartAndWait() error {
	err := server.Start()
	if err != nil {
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
