package server

import (
	"context"
	"fmt"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//Store represents the store in the cluster
type Store interface {
	Close() error
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

var emptyStoreCreateResponse = &pb.StoreCreateReply{}

func (server *server) createStoreAndStartRaft(name string, replicas int32) error {
	if _, ok := server.stores[name]; ok {
		return fmt.Errorf("store with name %s already exists", name)
	}
	// Verify if we have peers to achieve the replication
	if int32(len(server.peers)) < replicas-1 {
		return fmt.Errorf("cluster does not have required servers to achieve request " +
			"replication")
	}

	instance, err := NewRaftInstance(&RftArgs{})
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
	numOfReplicas := request.Replication
	err := server.createStoreAndStartRaft(request.GetName(), numOfReplicas)
	if err != nil {
		return emptyStoreCreateResponse, err
	}
	return emptyStoreCreateResponse, nil
}

//CreateStore Creates new store in the cluster
func (server *server) CreateStore(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	numOfReplicas := request.Replication
	err := server.createStoreAndStartRaft(request.GetName(), numOfReplicas)
	if err != nil {
		return emptyStoreCreateResponse, err
	}

	counter := 0
	for _, p := range server.peers {
		if counter >= int(numOfReplicas) {
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
		counter++
	}

	return emptyStoreCreateResponse, nil
}
