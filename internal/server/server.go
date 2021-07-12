package server

import (
	"context"
	"fmt"
	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

//YadosServer represents the YadosServer object in the cluster
type YadosServer struct {
	pb.UnimplementedYadosServiceServer

	self *pb.Member
	//OSSignalCh channel listening for the OS events
	OSSignalCh chan os.Signal
	peers      map[string]*pb.Member
	isTestMode bool
	logger     *logrus.Entry
	grpcServer *grpc.Server
}

// CreateNewServer Creates a new object of the YadosServer
func CreateNewServer(name string, address string, port int32) (*YadosServer, error) {
	server := YadosServer{
		self: &pb.Member{
			Port:    port,
			Address: address,
			Name:    name,
		},
		peers:      map[string]*pb.Member{},
		isTestMode: false,
		grpcServer: grpc.NewServer(),
	}

	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:[%YadosServer%] %time% - %msg% \n",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	server.logger = logger.WithFields(logrus.Fields{
		"YadosServer": server.self.Name,
	})
	server.OSSignalCh = make(chan os.Signal, 1)
	return &server, nil
}

func (server *YadosServer) startGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.self.Address, server.self.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pb.RegisterYadosServiceServer(server.grpcServer, server)
	go func() {
		server.grpcServer.Serve(lis)
	}()
	server.logger.Printf("server listening at %v", lis.Addr())
	return err
}

//StartAndWait start the server and wait for the OS signal
func (server *YadosServer) StartAndWait(peers []string) error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.startGrpcServer()
	if err != nil {
		return err
	}

	err = server.PostInit(peers)
	if err != nil {
		server.logger.Error(err)
		_ = server.StopServerFn()
		return err
	}

	<-server.OSSignalCh
	return nil
}

// HandleSignal Handles the CTRL-C signals
func (server *YadosServer) HandleSignal() {
	for {
		<-server.OSSignalCh
		log.Println("Exiting ... ")
		server.StopServerFn()
		os.Exit(0)
	}
}

// EnableTestMode To be ued from the test
func (server *YadosServer) EnableTestMode() {
	server.isTestMode = true
}

func (server *YadosServer) JoinWith(address string, port int32) error {
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

	_, err = peer.AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{Member: self})
	if err != nil {
		return err
	}

	return nil
}

// PostInit performs the post initialization
func (server *YadosServer) PostInit(peers []string) error {
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

//StopServerFn Stops the server
func (server *YadosServer) StopServerFn() error {
	server.logger.Info("Stopping the server ...")
	if !server.isTestMode {
		server.grpcServer.Stop()
		os.Exit(0)
	}
	return nil
}

// Status returns the status of the server
func (server *YadosServer) Status() error {
	return nil
}
