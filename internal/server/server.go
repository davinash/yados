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
	"syscall"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

//MemberServer represents the individual node in the cluster
type MemberServer struct {
	//Port on which server is running
	Port int32 `json:"port"`
	//Address ip address
	Address string `json:"address"`
	//Name unique name in the cluster
	Name string `json:"name"`
}

//YadosServer represents the YadosServer object in the cluster
type YadosServer struct {
	pb.UnimplementedYadosServiceServer

	self *MemberServer
	//OSSignalCh channel listening for the OS events
	OSSignalCh chan os.Signal
	peers      map[string]*MemberServer
	isTestMode bool
	logger     *logrus.Entry
	grpcServer *grpc.Server
}

// CreateNewServer Creates a new object of the YadosServer
func CreateNewServer(name string, address string, port int32) (*YadosServer, error) {
	server := YadosServer{
		self: &MemberServer{
			Port:    port,
			Address: address,
			Name:    name,
		},
		peers:      map[string]*MemberServer{},
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
func (server *YadosServer) StartAndWait(withPeer bool, peerAddress string, peerPort int32) error {
	signal.Notify(server.OSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go server.HandleSignal()

	err := server.startGrpcServer()
	if err != nil {
		return err
	}

	err = server.PostInit(withPeer, peerAddress, peerPort)
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

// PostInit performs the post initialization
func (server *YadosServer) PostInit(withPeer bool, peerAddress string, peerPort int32) error {
	server.logger.Info("Performing Post Initialization ...")

	if withPeer {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", peerAddress, peerPort), grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()
		peer := pb.NewYadosServiceClient(conn)
		_, err = peer.AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{
			Member: &pb.Member{
				Name:    server.self.Name,
				Address: server.self.Address,
				Port:    server.self.Port,
			},
		},
		)
		if err != nil {
			return err
		}

		//_, err := server.AddNewMemberInCluster(context.Background(), &pb.NewMemberRequest{
		//	Member: &pb.Member{
		//		Name:    "",
		//		Address: peerAddress,
		//		Port:    peerPort,
		//	},
		//})
		//if err != nil {
		//	return err
		//}

		//_, err := SendMessage(&MemberServer{
		//	Port:    peerPort,
		//	Address: peerAddress,
		//}, &Request{
		//	ID: AddNewMemberInCluster,
		//	Arguments: MemberServer{
		//		Port:    server.self.Port,
		//		Address: server.self.Address,
		//		Name:    server.self.Name,
		//	},
		//}, server.logger)
		//
		//if err != nil {
		//	server.logger.Error(err)
		//	return err
		//}
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
