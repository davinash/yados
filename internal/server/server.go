package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type MemberServer struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Server struct {
	self       *MemberServer
	listener   net.Listener
	OSSignalCh chan os.Signal
	peers      map[string]*MemberServer
	isTestMode bool
	logger     *logrus.Entry
}

func CreateNewServer(name string, address string, port int) (*Server, error) {
	server := Server{
		self: &MemberServer{
			Port:    port,
			Address: address,
			Name:    name,
		},
		peers:      map[string]*MemberServer{},
		isTestMode: false,
	}

	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:[%Server%] %time% - %msg% \n",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	server.logger = logger.WithFields(logrus.Fields{
		"Server": server.self.Name,
	})
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

func (server *Server) PostInit(withPeer bool, peerAddress string, peerPort int) error {
	server.logger.Info("Performing Post Initialization ...")

	if withPeer {
		_, err := SendMessage(&MemberServer{
			Port:    peerPort,
			Address: peerAddress,
		}, &Request{
			Id: AddNewMember,
			Arguments: MemberServer{
				Port:    server.self.Port,
				Address: server.self.Address,
				Name:    server.self.Name,
			},
		}, server.logger)

		if err != nil {
			server.logger.Error(err)
			return err
		}
		//if resp.Err != nil {
		//	server.logger.Error(resp.Err)
		//}
	}

	return nil
}

func (server *Server) StartAndWait(withPeer bool, peerAddress string, peerPort int) error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.startHttpServer()
	if err != nil {
		return err
	}
	err = server.PostInit(withPeer, peerAddress, peerPort)
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
