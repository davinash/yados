package server

import (
	"errors"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//ErrorInvalidFormat error for invalid format
var ErrorInvalidFormat = errors.New("invalid format for peers, use <ip-address>:port")

//Server Server interface
type Server interface {
	pb.YadosServiceServer
	Name() string
	Stop() error
	Address() string
	Port() int32
	Peers() []*pb.Peer
	Logger() *logrus.Entry
	SetLogLevel(level string)
	Serve([]*pb.Peer) error
	RPCServer() RPCServer
	Raft() Raft
	Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error)
	Self() *pb.Peer
	State() RaftState
}

type server struct {
	pb.UnimplementedYadosServiceServer
	mutex     sync.Mutex
	raft      Raft
	ready     <-chan interface{}
	quit      chan interface{}
	rpcServer RPCServer
	self      *pb.Peer
	logger    *logrus.Entry
}

//NewServer creates new instance of a server
func NewServer(name string, address string, port int32, loglevel string, ready <-chan interface{}) (Server, error) {
	srv := &server{}
	srv.self = &pb.Peer{
		Name:    name,
		Address: address,
		Port:    port,
	}
	srv.ready = ready

	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:%time% [%server%] %msg% \n",
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}
	srv.logger = logger.WithFields(logrus.Fields{
		"server": srv.self.Name,
	})
	srv.SetLogLevel(loglevel)
	srv.quit = make(chan interface{})

	return srv, nil
}

func (srv *server) Serve(peers []*pb.Peer) error {
	srv.logger.Infof("Starting Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.rpcServer = NewRPCServer(srv)
	err := srv.rpcServer.Start()
	if err != nil {
		srv.logger.Fatalf("failed to start the rpc, error = %v", err)
	}

	srv.raft, err = NewRaft(srv, peers, srv.ready)
	if err != nil {
		return err
	}
	return nil
}

func (srv *server) Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error) {
	reply, err := srv.RPCServer().Send(peer, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (srv *server) Name() string {
	return srv.self.Name
}

func (srv *server) Address() string {
	return srv.self.Address
}

func (srv *server) Port() int32 {
	return srv.self.Port
}

func (srv *server) Peers() []*pb.Peer {
	return srv.Raft().Peers()
}

func (srv *server) Logger() *logrus.Entry {
	return srv.logger
}

func (srv *server) RPCServer() RPCServer {
	return srv.rpcServer
}

func (srv *server) Raft() Raft {
	return srv.raft
}

func (srv *server) Self() *pb.Peer {
	return srv.self
}

func (srv *server) State() RaftState {
	return srv.Raft().State()
}

func (srv *server) SetLogLevel(level string) {
	switch level {
	case "trace":
		srv.logger.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		srv.logger.Logger.SetLevel(logrus.DebugLevel)
	case "info":
		srv.logger.Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		srv.logger.Logger.SetLevel(logrus.WarnLevel)
	case "error":
		srv.logger.Logger.SetLevel(logrus.ErrorLevel)
	default:
		srv.logger.Logger.SetLevel(logrus.InfoLevel)
	}
}

func (srv *server) Stop() error {
	srv.logger.Infof("Stopping Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	srv.Raft().Stop()
	srv.RPCServer().Stop()
	return nil
}
