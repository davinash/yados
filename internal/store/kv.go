package store

import (
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
)

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
func NewStore(name string) KVStore {
	s := &store{
		name:      name,
		kv:        make(map[string]string),
		storeType: pb.StoreType_Memory,
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

func (s *store) Close() error {
	return nil
}

func (s *store) Create(request *pb.StoreCreateRequest) {
	panic("implement me")
}

func (s *store) Delete(request *pb.StoreDeleteRequest) {
	panic("implement me")
}
