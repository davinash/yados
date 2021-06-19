package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Server struct {
	name       string
	port       int
	listener   net.Listener
	OSSignalCh chan os.Signal
}

func GetExistingServer() (*Server, error) {
	server := Server{}
	return &server, nil
}

func CreateNewServer(name string, port string) (*Server, error) {
	server := Server{
		name: name,
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	server.port = p
	server.OSSignalCh = make(chan os.Signal, 1)
	return &server, nil
}

func (server *Server) Start() error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	server.listener = listener

	go func() {
		http.Serve(server.listener, setupRouter(server))
	}()

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
