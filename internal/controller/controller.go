package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/davinash/yados/internal/rpc"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"google.golang.org/grpc"
)

//Controller struct for controller
type Controller struct {
	pb.UnimplementedControllerServiceServer

	address    string
	port       int32
	grpcServer *grpc.Server
	logger     *logrus.Logger

	members map[string]*pb.Peer
	mutex   sync.Mutex
}

//NewController create new instance of controller
func NewController(address string, port int32, loglevel string) *Controller {
	controller := &Controller{
		grpcServer: grpc.NewServer(),
		address:    address,
		port:       port,
		members:    make(map[string]*pb.Peer),
	}
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:%time% [Controller] %msg% \n",
			TimestampFormat: "Jan _2 15:04:05.000000000",
		},
	}
	controller.logger = logger
	controller.SetLogLevel(loglevel)
	return controller
}

//Address returns the address of this controller
func (c *Controller) Address() string {
	return c.address
}

//Port returns the port on which this controller is running
func (c *Controller) Port() int32 {
	return c.port
}

//SetLogLevel adds the log level
func (c *Controller) SetLogLevel(loglevel string) {
	switch loglevel {
	case "trace":
		c.logger.SetLevel(logrus.TraceLevel)
	case "debug":
		c.logger.SetLevel(logrus.DebugLevel)
	case "info":
		c.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		c.logger.SetLevel(logrus.WarnLevel)
	case "error":
		c.logger.SetLevel(logrus.ErrorLevel)
	default:
		c.logger.SetLevel(logrus.InfoLevel)
	}
}

//Start method to start the controller
func (c *Controller) Start() {
	c.logger.Infof("Starting a controller on [%s:%d]", c.address, c.port)
	pb.RegisterControllerServiceServer(c.grpcServer, c)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))
	if err != nil {
		panic(err)
	}
	go func() {
		err := c.grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
}

//Stop method to stop the controller
func (c *Controller) Stop() error {
	c.grpcServer.Stop()
	return nil
}

//Register perform the registration activity
func (c *Controller) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterReply, error) {
	c.logger.Debugf("Received Registration requirest for [%s:%s:%d]", request.Server.Name,
		request.Server.Address, request.Server.Port)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.members[request.Server.Name]; ok {
		return &pb.RegisterReply{}, fmt.Errorf("member with name %s already exists in the cluster",
			request.Server.Name)
	}
	reply := &pb.RegisterReply{}
	c.members[request.Server.Name] = request.Server

	for _, member := range c.members {
		peerConn, rpcClient, err := rpc.GetPeerConn(member.Address, member.Port)
		if err != nil {
			return reply, err
		}

		request := &pb.AddPeersRequest{
			Peers: make([]*pb.Peer, 0),
		}

		for _, m := range c.members {
			request.Peers = append(request.Peers, &pb.Peer{
				Name:    m.Name,
				Address: m.Address,
				Port:    m.Port,
			})
		}

		_, err = rpcClient.AddPeers(ctx, request)
		if err != nil {
			return nil, err
		}

		func(peerConn *grpc.ClientConn) {
			err := peerConn.Close()
			if err != nil {
				log.Printf("failed to close the connection, error = %v\n", err)
			}
		}(peerConn)
	}

	return reply, nil
}
