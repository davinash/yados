package store

import (
	"sync"

	"github.com/sirupsen/logrus"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//Store Interface for store
type Store interface {
	Name() string
	Delete() error
	Close() error
	Type() pb.StoreType
}

//KVStore interface for Key Value store
type KVStore interface {
	Store
	Put(*pb.PutRequest) error
	Get(*pb.GetRequest) string
	InternalMap() map[string]string
}

//SQLStore interface for Sql store
type SQLStore interface {
	Store
	ExecuteDDLQuery(*pb.DDLQueryRequest) (*pb.DDLQueryReply, error)
	ExecuteDMLQuery(*pb.DMLQueryRequest) (*pb.DMLQueryReply, error)
}

//Args arguments for creating new store
type Args struct {
	Name      string
	PLogDir   string
	Logger    *logrus.Entry
	StoreType pb.StoreType
}

type store struct {
	mutex     sync.RWMutex
	name      string
	kv        map[string]string
	storeType pb.StoreType
}

func (s *store) Type() pb.StoreType {
	return s.storeType
}

//NewStore creates a new store
func NewStore(args *Args) KVStore {
	s := &store{
		name:      args.Name,
		kv:        make(map[string]string),
		storeType: args.StoreType,
	}
	return s
}

func (s *store) Name() string {
	return s.name
}

func (s *store) InternalMap() map[string]string {
	return s.kv
}

func (s *store) Put(putRequest *pb.PutRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.kv[putRequest.Key] = putRequest.Value
	return nil
}

func (s *store) Get(getRequest *pb.GetRequest) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.kv[getRequest.Key]
}

func (s *store) Delete() error {
	panic("implement me")
}

func (s *store) Close() error {
	return nil
}
