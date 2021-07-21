package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

//Server server interface
type Server interface {
	pb.YadosServiceServer
	//SetOrCreateDataDir TODO Add later
	SetOrCreateDataDir(dir string) error
	//SetLogLevel TODO Add later
	SetLogLevel(level string)
	//Start TODO Add later
	Start(peers []string) error
	//Stop TODO Add later
	Stop() error
	//HandleSignal TODO Add later
	HandleSignal()
	//JoinWith TODO Add later
	JoinWith(address string, port int32) error
	//SetHealthCheckDuration TODO Add later
	SetHealthCheckDuration(duration int)
	//SetHealthCheckMaxWaitTime TODO Add later
	SetHealthCheckMaxWaitTime(maxTime int64)
	//Status TODO Add later
	Status() string
	//DataDirectory TODO Add later
	DataDirectory() string
	//Address TODO Add later
	Address() string
	//Port TODO Add later
	Port() int32
	//OSSignalChan TODO Add later
	OSSignalChan() chan os.Signal
	//Self TODO Add later
	Self() *pb.Member
	//PostInit TODO Add later
	PostInit(peers []string) error
	//Stores TODO Add later
	Stores() map[string]Store
	//EnableTestMode TODO Add later
	EnableTestMode()
}

//server represents the server object in the cluster
type server struct {
	pb.UnimplementedYadosServiceServer

	self *pb.Member
	//osSignalCh channel listening for the OS events
	osSignalCh        chan os.Signal
	peers             map[string]*pb.Member
	isTestMode        bool
	logger            *logrus.Entry
	grpcServer        *grpc.Server
	healthCheckChan   chan struct{}
	hcTriggerDuration int
	wg                sync.WaitGroup
	hcMaxWaitTime     int64
	mutex             *sync.RWMutex
	storeMutex        *sync.RWMutex
	stores            map[string]Store
	dataDirectory     string
}

//SrvArgs represents object required to start the server
type SrvArgs struct {
	Name    string
	Address string
	Port    int32
	HomeDir string
}

// CreateNewServer Creates a new object of the server
func CreateNewServer(args *SrvArgs) (Server, error) {
	server := server{
		self: &pb.Member{
			Port:    args.Port,
			Address: args.Address,
			Name:    args.Name,
		},
		peers:             map[string]*pb.Member{},
		isTestMode:        false,
		grpcServer:        grpc.NewServer(),
		healthCheckChan:   make(chan struct{}, 1),
		hcTriggerDuration: 30,
		hcMaxWaitTime:     180,
		mutex:             &sync.RWMutex{},
		storeMutex:        &sync.RWMutex{},
		stores:            make(map[string]Store),
	}

	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:[%server%] %time% - %msg% \n",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	server.logger = logger.WithFields(logrus.Fields{
		"server": server.self.Name,
	})

	err := server.SetOrCreateDataDir(args.HomeDir)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", args.Address, args.Port))
	if err != nil {
		return nil, fmt.Errorf("[CreateNewServer] failed to listen: %v", err)
	}
	lis.Close()

	server.osSignalCh = make(chan os.Signal, 1)
	return &server, nil
}

func (server *server) Self() *pb.Member {
	return server.self
}

func (server *server) Stores() map[string]Store {
	return server.stores
}

func (server *server) DataDirectory() string {
	return server.dataDirectory
}

func (server *server) Address() string {
	return server.self.Address
}

func (server *server) Port() int32 {
	return server.self.Port
}

func (server *server) OSSignalChan() chan os.Signal {
	return server.osSignalCh
}

//SetOrCreateDataDir Sets and create data directory
func (server *server) SetOrCreateDataDir(dir string) error {
	server.dataDirectory = dir
	if _, err := os.Stat(filepath.Join(dir, ".yados")); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(dir, ".yados"), os.ModePerm)
		if err != nil {
			return err
		}
	}
	server.logger.Infof("Data directory %s", server.DataDirectory())
	return nil
}

//SetLogLevel set the log level
func (server *server) SetLogLevel(level string) {
	switch level {
	case "trace":
		server.logger.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		server.logger.Logger.SetLevel(logrus.DebugLevel)
	case "info":
		server.logger.Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		server.logger.Logger.SetLevel(logrus.WarnLevel)
	case "error":
		server.logger.Logger.SetLevel(logrus.ErrorLevel)
	default:
		server.logger.Logger.SetLevel(logrus.InfoLevel)
	}
}

// StartGrpcServer Starts the GRPC server
func (server *server) StartGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.self.Address, server.self.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	pb.RegisterYadosServiceServer(server.grpcServer, server)
	go func() {
		server.grpcServer.Serve(lis)
	}()
	server.logger.Infof("server listening at %v", lis.Addr())
	return err
}

//Start start the server and wait for the OS signal
func (server *server) Start(peers []string) error {
	signal.Notify(server.osSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.StartGrpcServer()
	if err != nil {
		return err
	}
	err = server.PostInit(peers)
	if err != nil {
		server.logger.Error(err)
		_ = server.Stop()
		return err
	}
	err = server.StartHealthCheck()
	if err != nil {
		server.logger.Error(err)
		_ = server.Stop()
		return err
	}
	return nil
}

// HandleSignal Handles the CTRL-C signals
func (server *server) HandleSignal() {
	for {
		<-server.osSignalCh
		err := server.Stop()
		if err != nil {
			return
		}
		if !server.isTestMode {
			os.Exit(0)
		}
	}
}

// EnableTestMode To be ued from the test
func (server *server) EnableTestMode() {
	server.isTestMode = true
}

// JoinWith admits the new members in the existing cluster
func (server *server) JoinWith(address string, port int32) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	conn, peer, err := GetPeerConn(address, port)
	if err != nil {
		return err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			server.logger.Warnf("Failed to close connection, Error = %v", err)
		}
	}(conn)

	self := &pb.Member{
		Name:    server.self.Name,
		Address: server.self.Address,
		Port:    server.self.Port,
	}

	remotePeer, err := peer.AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{Member: self})
	if err != nil {
		return err
	}
	// Add the remote to self peer list
	server.peers[remotePeer.Member.Name] = remotePeer.Member
	return nil
}

// PostInit performs the post initialization
func (server *server) PostInit(peers []string) error {
	server.logger.Info("Performing Post Initialization ...")
	if len(peers) == 0 {
		return nil
	}
	for _, p := range peers {
		split := strings.Split(p, ":")
		if len(split) != 2 {
			return fmt.Errorf("invalid format for peers, use <ip-address>:port")
		}
		port, err := strconv.Atoi(split[1])
		if err != nil {
			return fmt.Errorf("invalid format for peers, use <ip-address>:port")
		}
		err = server.JoinWith(split[0], int32(port))
		if err != nil {
			return err
		}
	}
	return nil
}

// SetHealthCheckDuration Sets the time duration to run the health check
func (server *server) SetHealthCheckDuration(duration int) {
	server.hcTriggerDuration = duration
}

//SetHealthCheckMaxWaitTime Set the max time to wait for members response before declaring
//it is un-healthy
func (server *server) SetHealthCheckMaxWaitTime(maxTime int64) {
	server.hcMaxWaitTime = maxTime
}

//Stop Stops the server
func (server *server) Stop() error {
	server.logger.Infof("Stopping the server [%s:%d] ...", server.self.Address, server.self.Port)
	server.StopHealthCheck()

	for _, s := range server.stores {
		s.Close()
	}

	//server.raftInstance.Stop()
	server.grpcServer.Stop()

	server.wg.Wait()
	if !server.isTestMode {
		os.Exit(0)
	}
	return nil
}

// Status returns the status of the server
func (server *server) Status() string {
	return ""
}
