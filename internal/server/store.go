package server

//Store Interface for store
type Store interface {
	Open() error
	Close() error
	Put() error
	Get() error
}

//StoreArgs arguments for creating new store
type StoreArgs struct {
}

//NewStore creates a new store
func NewStore(args *StoreArgs) Store {
	s := &store{}
	return s
}

type store struct {
}

func (s *store) Open() error {
	panic("implement me")
}

func (s *store) Close() error {
	panic("implement me")
}

func (s *store) Put() error {
	panic("implement me")
}

func (s *store) Get() error {
	panic("implement me")
}
