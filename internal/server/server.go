package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/davinash/yados/internal/store"

	"github.com/davinash/yados/internal/events"
	"github.com/davinash/yados/internal/raft"
	"github.com/davinash/yados/internal/rpc"
	"github.com/davinash/yados/internal/wal"
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
	Logger() *logrus.Logger
	SetLogLevel(level string)
	Serve() error
	RPCServer() rpc.Server
	Raft() raft.Raft
	Self() *pb.Peer
	State() raft.State
	WALDir() string
	WAL() wal.Wal
	HTTPPort() int
	StartHTTPServer()
	StopHTTPServer() error
	StoreManager() store.Manager
	Controller() []controller

	// EventHandler for test purpose
	EventHandler() *events.Events
}

type controller struct {
	address string
	port    int32
}

type server struct {
	pb.UnimplementedYadosServiceServer
	mutex     sync.Mutex
	raft      raft.Raft
	quit      chan interface{}
	rpcServer rpc.Server
	self      *pb.Peer
	logger    *logrus.Logger
	walDir    string
	wal       wal.Wal
	ev        *events.Events

	httpPort       int
	httpServerWait sync.WaitGroup
	httpSrv        *http.Server
	storageMgr     store.Manager

	controllers []controller
	isTestMode  bool
}

//NewServerArgs argument structure for new server
type NewServerArgs struct {
	Name        string
	Address     string
	Port        int32
	Loglevel    string
	WalDir      string
	IsTestMode  bool
	HTTPPort    int
	Controllers []string
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
	srv.walDir = args.WalDir

	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &easy.Formatter{
			LogFormat:       "[%lvl%]:%time% %msg% \n",
			TimestampFormat: "Jan _2 15:04:05.000000000",
		},
	}

	srv.logger = logger
	srv.SetLogLevel(args.Loglevel)
	srv.quit = make(chan interface{})

	srv.httpPort = args.HTTPPort

	if args.IsTestMode {
		srv.isTestMode = true
		srv.ev = events.NewEvents()
	}

	srv.controllers = make([]controller, 0)
	for _, c := range args.Controllers {
		split := strings.Split(c, ":")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid format for peers, use <ip-address:port>")
		}
		port, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, fmt.Errorf("invalid format for peers, use <ip-address:port>")
		}
		srv.controllers = append(srv.controllers, controller{
			address: split[0],
			port:    int32(port),
		})
	}
	return srv, nil
}

func (srv *server) Serve() error {
	srv.GetOrCreateWAL()
	srv.rpcServer = rpc.NewRPCServer(srv.Name(), srv.self.Address, srv.self.Port, srv.logger)
	pb.RegisterYadosServiceServer(srv.rpcServer.GrpcServer(), srv)
	srv.rpcServer.Start()

	srv.storageMgr = store.NewStoreManger(srv.logger, srv.walDir)

	raft, err := raft.NewRaft(&raft.Args{
		IsTestMode:   srv.isTestMode,
		Logger:       srv.logger,
		PstLog:       srv.wal,
		RPCServer:    srv.rpcServer,
		Server:       srv.Self(),
		EventHandler: srv.EventHandler(),
		StorageMgr:   srv.StoreManager(),
	})
	if err != nil {
		panic(err)
	}
	srv.raft = raft

	for _, controller := range srv.Controller() {
		srv.logger.Debugf("[%s] Registration with controller", srv.Name())
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", controller.address, controller.port), grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("[%s] failed to connect with controller[%s:%d], "+
				"error = %w", srv.Name(), controller.address, controller.port, err)
		}
		controller := pb.NewControllerServiceClient(conn)
		_, err = controller.Register(context.Background(), &pb.RegisterRequest{Server: srv.self})
		if err != nil {
			return err
		}
		err = conn.Close()
		if err != nil {
			srv.logger.Warnf("[%s] Failed to close the connection, Error = %v", srv.Name(), err)
		}
	}

	srv.StartHTTPServer()
	srv.raft.Start()
	srv.logger.Infof("Started Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())

	return nil
}

func (srv *server) GetOrCreateWAL() {
	d := filepath.Join(srv.walDir, srv.Name(), "wal")
	if err := os.MkdirAll(d, os.ModePerm); err != nil {
		panic(err)
	}
	srv.walDir = d
	s := wal.NewWAL(srv.walDir, srv.logger, srv.EventHandler(), srv.isTestMode)

	srv.wal = s
	srv.wal.Open()
}

func (srv *server) Controller() []controller {
	return srv.controllers
}

func (srv *server) WALDir() string {
	return srv.walDir
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

func (srv *server) Logger() *logrus.Logger {
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

func (srv *server) WAL() wal.Wal {
	return srv.wal
}

func (srv *server) EventHandler() *events.Events {
	return srv.ev
}

func (srv *server) HTTPPort() int {
	return srv.httpPort
}

func (srv *server) StoreManager() store.Manager {
	return srv.storageMgr
}

func (srv *server) SetLogLevel(level string) {
	switch level {
	case "trace":
		srv.logger.SetLevel(logrus.TraceLevel)
	case "debug":
		srv.logger.SetLevel(logrus.DebugLevel)
	case "info":
		srv.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		srv.logger.SetLevel(logrus.WarnLevel)
	case "error":
		srv.logger.SetLevel(logrus.ErrorLevel)
	default:
		srv.logger.SetLevel(logrus.InfoLevel)
	}
}

func (srv *server) Stop() error {
	srv.logger.Infof("Stopping Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())

	err := srv.StopHTTPServer()
	if err != nil {
		return err
	}

	srv.Raft().Stop()
	err = srv.RPCServer().Stop()
	if err != nil {
		srv.logger.Errorf("failed to stop grpc server, Error = %v", err)
		return err
	}

	for _, controller := range srv.Controller() {
		srv.logger.Debugf("[%s] UnRegistration with controller", srv.Name())
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", controller.address, controller.port), grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("[%s] failed to connect with controller[%s:%d], "+
				"error = %w", srv.Name(), controller.address, controller.port, err)
		}
		controller := pb.NewControllerServiceClient(conn)
		_, err = controller.UnRegister(context.Background(), &pb.UnRegisterRequest{Server: srv.self})
		if err != nil {
			return err
		}
		err = conn.Close()
		if err != nil {
			srv.logger.Warnf("[%s] Failed to close the connection, Error = %v", srv.Name(), err)
		}
	}
	time.Sleep(10 * time.Millisecond)

	err = srv.StoreManager().Close()
	if err != nil {
		return err
	}

	err = srv.WAL().Close()
	if err != nil {
		srv.logger.Errorf("failed to close the wal, Error = %v", err)
		return err
	}
	srv.httpServerWait.Wait()
	srv.logger.Infof("Stopped Server %s on [%s:%d]", srv.Name(), srv.Address(), srv.Port())
	return nil
}

func (srv *server) StartHTTPServer() {
	if srv.HTTPPort() == -1 {
		return
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
			panic(err)
		}
	}()
	srv.logger.Infof("Started HTTP server on %s:%d", srv.Self().Address, srv.HTTPPort())
}

func (srv *server) StopHTTPServer() error {
	if srv.httpSrv != nil {
		if err := srv.httpSrv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}
	return nil
}
