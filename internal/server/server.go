package server

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/protobuf/proto"

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
	PLog() PLog
	StoreCreate(request *pb.StoreCreateRequest) error
	Apply(entry *pb.LogEntry) error
	Stores() map[string]Store

	// SetEventHandler For Test Purpose
	SetEventHandler(Events)
	// EventHandler for test purpose
	EventHandler() Events
}

type server struct {
	pb.UnimplementedYadosServiceServer
	mutex      sync.Mutex
	raft       Raft
	quit       chan interface{}
	rpcServer  RPCServer
	self       *pb.Peer
	logger     *logrus.Entry
	logDir     string
	pLog       PLog
	stores     map[string]Store
	isTestMode bool
	ev         Events
}

//NewServerArgs argument structure for new server
type NewServerArgs struct {
	Name       string
	Address    string
	Port       int32
	Loglevel   string
	LogDir     string
	IsTestMode bool
}

//NewServer creates new instance of a server
func NewServer(args *NewServerArgs) (Server, error) {
	srv := &server{
		isTestMode: false,
	}
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
	srv.stores = make(map[string]Store)

	if args.IsTestMode {
		srv.isTestMode = true
		srv.ev = NewEvents()
	}

	return srv, nil
}

func (srv *server) GetOrCreateStorage() error {
	d := filepath.Join(srv.logDir, srv.Name(), "log")
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	srv.logDir = d

	store, err := NewPLog(srv.LogDir(), srv.Logger())
	if err != nil {
		return err
	}
	srv.pLog = store
	err = srv.pLog.Open()
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

	args := RaftArgs{
		srv:        srv,
		peers:      peers,
		isTestMode: srv.isTestMode,
	}
	srv.raft, err = NewRaft(&args)
	if err != nil {
		return err
	}
	srv.raft.Start()

	return nil
}

func (srv *server) StoreCreate(request *pb.StoreCreateRequest) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	s := NewStore(&StoreArgs{})
	srv.Stores()[request.Name] = s

	return nil
}

func (srv *server) Apply(entry *pb.LogEntry) error {
	switch entry.CmdType {
	case pb.CommandType_CreateStore:
		var scr pb.StoreCreateRequest
		err := proto.Unmarshal(entry.Command, &scr)
		if err != nil {
			return err
		}
		err = srv.StoreCreate(&scr)
		if err != nil {
			return err
		}
	case pb.CommandType_Put:
		var putRequest pb.PutRequest
		err := proto.Unmarshal(entry.Command, &putRequest)
		if err != nil {
			return err
		}
		return srv.Stores()[putRequest.StoreName].Put(&putRequest)
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

func (srv *server) PLog() PLog {
	return srv.pLog
}

func (srv *server) Stores() map[string]Store {
	return srv.stores
}

func (srv *server) SetEventHandler(ev Events) {
	srv.ev = ev
}
func (srv *server) EventHandler() Events {
	return srv.ev
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
	err = srv.PLog().Close()
	if err != nil {
		srv.logger.Errorf("failed to close the pLog, Error = %v", err)
		return err
	}
	srv.logger.Infof("Stopped Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	return nil
}
