package store

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

//Store Interface for store
type Store interface {
	Name() string
	Close() error
	Create(*pb.StoreCreateRequest)
	Delete(*pb.StoreDeleteRequest)
}

//SQLStore interface for Sql store
type SQLStore interface {
	Store
	Execute(*pb.ExecuteQueryRequest) (*pb.ExecuteQueryReply, error)
	Query(*pb.QueryRequest) (*pb.QueryReply, error)
}

//Args arguments for creating new store
type Args struct {
	Name   string
	WALDir string
	Logger *logrus.Logger
}

//Manager interface for Storage management
type Manager interface {
	Create(*pb.StoreCreateRequest) error
	Delete(*pb.StoreDeleteRequest) error
	Apply(entry *pb.WalEntry) error
	Close() error
	Stores() map[string]Store
}

type manager struct {
	mutex  sync.RWMutex
	stores map[string]Store
	logger *logrus.Logger
	walDir string
}

//NewStoreManger new instance of storage manager
func NewStoreManger(logger *logrus.Logger, walDir string) Manager {
	m := &manager{
		logger: logger,
		walDir: walDir,
	}
	m.stores = make(map[string]Store)
	return m
}

func (sm *manager) Create(request *pb.StoreCreateRequest) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	s, err := NewSqliteStore(&Args{
		Name:   request.Name,
		WALDir: sm.walDir,
		Logger: sm.logger,
	})
	if err != nil {
		return err
	}
	sm.Stores()[request.Name] = s
	return nil
}

func (sm *manager) Close() error {
	for _, store := range sm.Stores() {
		err := store.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm *manager) Delete(request *pb.StoreDeleteRequest) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.Stores(), request.StoreName)

	return nil
}

func (sm *manager) Stores() map[string]Store {
	return sm.stores
}

func (sm *manager) Apply(entry *pb.WalEntry) error {
	switch entry.CmdType {

	case pb.CommandType_CreateStore:
		var req pb.StoreCreateRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		err = sm.Create(&req)
		if err != nil {
			return err
		}

	case pb.CommandType_DeleteStore:
		var req pb.StoreDeleteRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		err = sm.Delete(&req)
		if err != nil {
			return err
		}

	case pb.CommandType_SqlDDL:
		var req pb.ExecuteQueryRequest
		err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
		if err != nil {
			return err
		}
		q, err := (sm.Stores()[req.StoreName].(SQLStore)).Execute(&req)
		if err != nil {
			return err
		}
		sm.logger.Debug(q)
	}
	return nil
}
