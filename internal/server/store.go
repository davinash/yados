package server

import (
	"context"
	"fmt"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//Store represents the store in the cluster
type Store interface {
	Close() error
	Name() string
	RaftInstance() Raft
}

type store struct {
	name string
	raft Raft
}

//NewStore creates new store
func NewStore(name string, raft Raft) Store {
	return &store{
		name: name,
		raft: raft,
	}
}

func (s *store) Close() error {
	err := s.raft.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (s *store) Name() string {
	return s.name
}
func (s *store) RaftInstance() Raft {
	return s.raft
}

var emptyStoreCreateResponse = &pb.StoreCreateReply{}

func (server *server) createStoreAndStartRaft(name string) error {
	if _, ok := server.stores[name]; ok {
		return fmt.Errorf("store with name %s already exists", name)
	}

	args := &RftArgs{
		ServerName: server.self.Name,
	}
	if server.isTestMode {
		args.IsTestMode = true
	}
	instance, err := NewRaftInstance(args)
	if err != nil {
		return err
	}
	instance.SetLogger(server.logger)
	instance.SetName(name)
	err = instance.Start()
	if err != nil {
		return err
	}
	newStore := NewStore(name, instance)
	server.stores[name] = newStore

	return nil
}

//CreateStoreSecondary creates the store on the secondary
func (server *server) CreateStoreSecondary(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	err := server.createStoreAndStartRaft(request.GetName())
	if err != nil {
		return emptyStoreCreateResponse, err
	}
	return emptyStoreCreateResponse, nil
}

//CreateStore Creates new store in the cluster
func (server *server) CreateStore(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	replicas := request.Replication

	// Verify if we have peers to achieve the replication
	if int32(len(server.peers)) < replicas {
		return emptyStoreCreateResponse, fmt.Errorf("cluster does not have required servers to achieve request " +
			"replication")
	}

	err := server.createStoreAndStartRaft(request.GetName())
	if err != nil {
		return emptyStoreCreateResponse, err
	}

	counter := 0
	for _, p := range server.peers {

		if counter >= int(replicas) {
			break
		}
		conn, p2, err := GetPeerConn(p.Address, p.Port)
		if err != nil {
			return emptyStoreCreateResponse, err
		}
		_, err = p2.CreateStoreSecondary(ctx, request)
		if err != nil {
			err = conn.Close()
			return emptyStoreCreateResponse, err
		}
		err = conn.Close()
		if err != nil {
			server.logger.Warnf("Failed to close the connection, error = %v", err)
		}

		server.Stores()[request.GetName()].RaftInstance().AddPeer(p)

		counter++
	}

	return emptyStoreCreateResponse, nil
}
