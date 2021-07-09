package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

//MemberServer represents the individual node in the cluster
type MemberServer struct {
	//Port on which server is running
	Port int `json:"port"`
	//Address ip address
	Address string `json:"address"`
	//Name unique name in the cluster
	Name string `json:"name"`
}

//Server represents the Server object in the cluster
type Server struct {
	self     *MemberServer
	listener net.Listener
	//OSSignalCh channel listening for the OS events
	OSSignalCh chan os.Signal
	peers      map[string]*MemberServer
	isTestMode bool
	logger     *logrus.Entry
}

// CreateNewServer Creates a new object of the Server
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

// HandleSignal Handles the CTRL-C signals
func (server *Server) HandleSignal() {
	for {
		<-server.OSSignalCh
		log.Println("Exiting ... ")
		server.StopServerFn()
		os.Exit(0)
	}
}

// EnableTestMode To be ued from the test
func (server *Server) EnableTestMode() {
	server.isTestMode = true
}

func (server *Server) startHTTPServer() error {
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

// PostInit performs the post initialization
func (server *Server) PostInit(withPeer bool, peerAddress string, peerPort int) error {
	server.logger.Info("Performing Post Initialization ...")

	if withPeer {
		_, err := SendMessage(&MemberServer{
			Port:    peerPort,
			Address: peerAddress,
		}, &Request{
			ID: AddNewMemberInCluster,
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
	}
	return nil
}

//StartAndWait start the server and wait for the OS signal
func (server *Server) StartAndWait(withPeer bool, peerAddress string, peerPort int) error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.startHTTPServer()
	if err != nil {
		return err
	}
	err = server.PostInit(withPeer, peerAddress, peerPort)
	if err != nil {
		server.StopServerFn()
		return err
	}
	<-server.OSSignalCh
	return nil
}

//StopServerFn Stops the server
func (server *Server) StopServerFn() error {
	server.logger.Info("Stopping the server ...")
	if !server.isTestMode {
		if err := server.listener.Close(); err != nil {
			log.Printf("failed to stop the http server, Error = %v\n", err)
			return err
		}
		os.Exit(0)
	}
	return nil
}

// Status returns the status of the server
func (server *Server) Status() error {
	return nil
}
