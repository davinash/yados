package server

import (
	"errors"
	"os"
	"path/filepath"
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
	Peers() map[string]*pb.Peer
	Logger() *logrus.Entry
	SetLogLevel(level string)
	Serve([]*pb.Peer) error
	RPCServer() RPCServer
	Raft() Raft
	Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error)
	Self() *pb.Peer
	State() RaftState
	LogDir() string
	Store() Store
}

type server struct {
	pb.UnimplementedYadosServiceServer
	mutex     sync.Mutex
	raft      Raft
	quit      chan interface{}
	rpcServer RPCServer
	self      *pb.Peer
	logger    *logrus.Entry
	logDir    string
	store     Store
}

//NewServerArgs argument structure for new server
type NewServerArgs struct {
	Name     string
	Address  string
	Port     int32
	Loglevel string
	LogDir   string
}

//NewServer creates new instance of a server
func NewServer(args *NewServerArgs) (Server, error) {
	srv := &server{}
	srv.self = &pb.Peer{
		Name:    args.Name,
		Address: args.Address,
		Port:    args.Port,
	}
	srv.logDir = args.LogDir

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

	srv.SetLogLevel(args.Loglevel)
	srv.quit = make(chan interface{})

	return srv, nil
}

func (srv *server) GetOrCreateStorage() error {
	d := filepath.Join(srv.logDir, srv.Name(), "log")
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	srv.logDir = d

	store, err := NewStorage(srv.LogDir(), srv.Logger())
	if err != nil {
		return err
	}
	srv.store = store
	err = srv.store.Open()
	if err != nil {
		return err
	}
	return nil
}

func (srv *server) Serve(peers []*pb.Peer) error {
	err := srv.GetOrCreateStorage()
	if err != nil {
		return err
	}
	srv.logger.Infof("Starting Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.rpcServer = NewRPCServer(srv)
	err = srv.rpcServer.Start()
	if err != nil {
		srv.logger.Fatalf("failed to start the rpc, error = %v", err)
	}

	srv.raft, err = NewRaft(srv, peers)
	if err != nil {
		return err
	}
	srv.raft.Start()

	return nil
}

func (srv *server) Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error) {
	reply, err := srv.RPCServer().Send(peer, serviceMethod, args)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (srv *server) LogDir() string {
	return srv.logDir
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

func (srv *server) Peers() map[string]*pb.Peer {
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

func (srv *server) Store() Store {
	return srv.store
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
	err := srv.RPCServer().Stop()
	if err != nil {
		srv.logger.Errorf("failed to stop grpc server, Error = %v", err)
		return err
	}
	err = srv.Store().Close()
	if err != nil {
		srv.logger.Errorf("failed to close the store, Error = %v", err)
		return err
	}
	srv.logger.Infof("Stopped Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	return nil
}
