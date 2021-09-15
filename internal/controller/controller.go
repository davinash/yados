package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

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
	leader  *pb.Peer
	wg      sync.WaitGroup
	quit    chan interface{}
}

//NewController create new instance of controller
func NewController(address string, port int32, loglevel string) *Controller {
	controller := &Controller{
		quit:       make(chan interface{}),
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

func (c *Controller) findLeader() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, m := range c.members {
		if c.IsLeader(m) {
			c.leader = m
			break
		}
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

	c.wg.Add(1)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		defer c.wg.Done()
		for {
			select {
			case <-c.quit:
				return
			case <-ticker.C:
				c.findLeader()
			}
		}
	}()
}

//Stop method to stop the controller
func (c *Controller) Stop() error {
	c.logger.Debugf("Stopping the controller")
	close(c.quit)
	c.grpcServer.Stop()
	c.wg.Wait()
	return nil
}

//IsLeader helper function to get the leader
func (c *Controller) IsLeader(member *pb.Peer) bool {
	peerConn, rpcClient, err := rpc.GetPeerConn(member.Address, member.Port)
	if err != nil {
		return false
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			c.logger.Warnf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)
	request := &pb.StatusRequest{}
	request.Id = uuid.New().String()
	reply, err := rpcClient.PeerStatus(context.Background(), request)
	if err != nil {
		return false
	}
	return reply.IsLeader
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

//GetLeader returns the leader in the cluster at this point of time
func (c *Controller) GetLeader(ctx context.Context, request *pb.GetLeaderRequest) (*pb.GetLeaderReply, error) {
	reply := &pb.GetLeaderReply{}
	c.mutex.Lock()
	reply.Leader = c.leader
	c.mutex.Unlock()
	return reply, nil
}

//UnRegister function to unregister the member from the cluster
func (c *Controller) UnRegister(ctx context.Context, request *pb.UnRegisterRequest) (*pb.UnRegisterReply, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, p := range c.members {
		if p.Name != request.Server.Name {
			peerConn, rpcClient, err := rpc.GetPeerConn(p.Address, p.Port)
			if err != nil {
				return &pb.UnRegisterReply{}, err
			}
			c.logger.Debugf("Removing peer [%s:%s:%d] from [%s:%s:%d]", request.Server.Name,
				request.Server.Address, request.Server.Port, p.Name, p.Address, p.Port)
			_, err = rpcClient.RemovePeer(ctx, &pb.RemovePeerRequest{
				Peer: request.Server,
				Id:   uuid.New().String(),
			})
			if err != nil {
				return &pb.UnRegisterReply{}, err
			}
			err = peerConn.Close()
			if err != nil {
				c.logger.Warnf("failed to close the connection with [%s:%s:%d], "+
					"Error= %v", p.Name, p.Address, p.Port, err)
			}
		}
	}

	if _, ok := c.members[request.Server.Name]; !ok {
		return &pb.UnRegisterReply{}, nil
	}
	c.leader = nil
	delete(c.members, request.Server.Name)

	return &pb.UnRegisterReply{}, nil
}
