package store

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//Store Interface for store
type Store interface {
	Put(*pb.PutRequest) error
	Get(*pb.GetRequest) string
	Name() string
	Delete() error
	InternalMap() map[string]string
}

//Args arguments for creating new store
type Args struct {
	name string
}

type store struct {
	mutex sync.RWMutex
	name  string
	kv    map[string]string
}

//NewStore creates a new store
func NewStore(args *Args) Store {
	s := &store{
		name: args.name,
		kv:   make(map[string]string),
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
