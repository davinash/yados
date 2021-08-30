package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	store_sqlite "github.com/davinash/yados/internal/sqlite"

	"github.com/davinash/yados/internal/events"
	"github.com/davinash/yados/internal/plog"
	"github.com/davinash/yados/internal/raft"
	"github.com/davinash/yados/internal/rpc"
	"github.com/davinash/yados/internal/store"
	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/anypb"

	"google.golang.org/protobuf/proto"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"

	pb "github.com/davinash/yados/internal/proto/gen"
)

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
	RPCServer() rpc.Server
	Raft() raft.Raft
	Send(peer *pb.Peer, serviceMethod string, args interface{}) (interface{}, error)
	Self() *pb.Peer
	State() raft.State
	LogDir() string
	PLog() plog.PLog
	StoreCreate(request *pb.StoreCreateRequest) error
	StoreDelete(request *pb.StoreDeleteRequest) error
	Apply(entry *pb.LogEntry) error
	Stores() map[string]store.Store

	HTTPPort() int
	StartHTTPServer() error
	StopHTTPServer() error

	// EventHandler for test purpose
	EventHandler() *events.Events
}

type server struct {
	pb.UnimplementedYadosServiceServer
	mutex      sync.Mutex
	raft       raft.Raft
	quit       chan interface{}
	rpcServer  rpc.Server
	self       *pb.Peer
	logger     *logrus.Entry
	logDir     string
	pLog       plog.PLog
	stores     map[string]store.Store
	isTestMode bool
	ev         *events.Events

	httpPort       int
	httpServerWait sync.WaitGroup
	httpSrv        *http.Server
}

//NewServerArgs argument structure for new server
type NewServerArgs struct {
	Name       string
	Address    string
	Port       int32
	Loglevel   string
	LogDir     string
	IsTestMode bool
	HTTPPort   int
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
	srv.stores = make(map[string]store.Store)

	srv.httpPort = args.HTTPPort

	if args.IsTestMode {
		srv.isTestMode = true
		srv.ev = events.NewEvents()
	}

	return srv, nil
}

func (srv *server) Serve(peers []*pb.Peer) error {
	err := srv.GetOrCreateStorage()
	if err != nil {
		return err
	}
	srv.logger.Infof("Starting Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	srv.rpcServer = rpc.NewRPCServer(srv.Name(), srv.self.Address, srv.self.Port, srv.logger)

	pb.RegisterYadosServiceServer(srv.rpcServer.GrpcServer(), srv)

	err = srv.rpcServer.Start()
	if err != nil {
		srv.logger.Fatalf("failed to start the rpc, error = %v", err)
	}

	args := raft.Args{
		IsTestMode:   srv.isTestMode,
		Apply:        srv.Apply,
		Logger:       srv.logger,
		PstLog:       srv.pLog,
		RPCServer:    srv.rpcServer,
		Server:       srv.Self(),
		EventHandler: srv.EventHandler(),
	}
	srv.raft, err = raft.NewRaft(&args)
	if err != nil {
		return err
	}

	for _, p := range peers {
		_, err := srv.Send(p, "server.AddNewMember", &pb.NewPeerRequest{
			Id: uuid.New().String(),
			NewPeer: &pb.Peer{
				Name:    srv.Name(),
				Address: srv.Address(),
				Port:    srv.Port(),
			},
		})
		if err != nil {
			srv.logger.Fatalf("failed add peer , err = %v", err)
		}
		// add this peer to self
		err = srv.raft.AddPeer(p)
		if err != nil {
			srv.logger.Fatalf("failed add peer , err = %v", err)
		}
	}
	srv.raft.Start()

	err = srv.StartHTTPServer()
	if err != nil {
		return err
	}

	return nil
}

func (srv *server) GetOrCreateStorage() error {
	d := filepath.Join(srv.logDir, srv.Name(), "log")
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	srv.logDir = d

	s, err1 := plog.NewPLog(srv.logDir, srv.logger, srv.EventHandler(), srv.isTestMode)
	if err1 != nil {
		return err1
	}
	srv.pLog = s
	err = srv.pLog.Open()
	if err != nil {
		return err
	}

	return nil
}

func (srv *server) StoreCreate(request *pb.StoreCreateRequest) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	if request.Type == pb.StoreType_Sqlite {
		s, err := store_sqlite.NewSqliteStore(&store.Args{
			Name:    request.Name,
			PLogDir: srv.LogDir(),
			Logger:  srv.Logger(),
		})
		if err != nil {
			return err
		}
		srv.Stores()[request.Name] = s
		return nil
	}
	s := store.NewStore(&store.Args{
		StoreType: pb.StoreType_Memory,
	})
	srv.Stores()[request.Name] = s

	return nil
}

func (srv *server) StoreDelete(request *pb.StoreDeleteRequest) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	delete(srv.Stores(), request.StoreName)

	return nil
}

func (srv *server) Apply(entry *pb.LogEntry) error {
	switch entry.CmdType {

	case pb.CommandType_CreateStore:
		var req pb.StoreCreateRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		err = srv.StoreCreate(&req)
		if err != nil {
			return err
		}

	case pb.CommandType_Put:
		var req pb.PutRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		return (srv.Stores()[req.StoreName].(store.KVStore)).Put(&req)

	case pb.CommandType_DeleteStore:
		var req pb.StoreDeleteRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		err = srv.StoreDelete(&req)
		if err != nil {
			return err
		}

	case pb.CommandType_SqlDDL:
		var req pb.DDLQueryRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		q, err := (srv.Stores()[req.StoreName].(store.SQLStore)).ExecuteDDLQuery(&req)
		if err != nil {
			return err
		}
		srv.logger.Debug(q)
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

func (srv *server) RPCServer() rpc.Server {
	return srv.rpcServer
}

func (srv *server) Raft() raft.Raft {
	return srv.raft
}

func (srv *server) Self() *pb.Peer {
	return srv.self
}

func (srv *server) State() raft.State {
	return srv.Raft().State()
}

func (srv *server) PLog() plog.PLog {
	return srv.pLog
}

func (srv *server) Stores() map[string]store.Store {
	return srv.stores
}

func (srv *server) EventHandler() *events.Events {
	return srv.ev
}

func (srv *server) HTTPPort() int {
	return srv.httpPort
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
	err := srv.StopHTTPServer()
	if err != nil {
		return err
	}
	// Close all the stores
	for _, v := range srv.Stores() {
		err := v.Close()
		if err != nil {
			return err
		}
	}

	srv.Raft().Stop()
	err = srv.RPCServer().Stop()
	if err != nil {
		srv.logger.Errorf("failed to stop grpc server, Error = %v", err)
		return err
	}
	err = srv.PLog().Close()
	if err != nil {
		srv.logger.Errorf("failed to close the pLog, Error = %v", err)
		return err
	}
	srv.httpServerWait.Wait()
	srv.logger.Infof("Stopped Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	return nil
}

func (srv *server) StartHTTPServer() error {
	if srv.HTTPPort() == -1 {
		return nil
	}
	httpHandler := NewHTTPHandler(srv)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	srv.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", srv.Self().Address, srv.HTTPPort()),
		Handler: httpHandler.Router(),
	}
	srv.httpServerWait.Add(1)
	go func() {
		defer srv.httpServerWait.Done()
		if err := srv.httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. Port in use?
			srv.logger.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	srv.logger.Infof("Started HTTP server on %s:%d", srv.Self().Address, srv.HTTPPort())
	return nil
}

func (srv *server) StopHTTPServer() error {
	if srv.httpSrv != nil {
		if err := srv.httpSrv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}
	return nil
}
