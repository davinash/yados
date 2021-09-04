package server

import (
	"context"
	"fmt"
	"log"

	"github.com/davinash/yados/internal/rpc"

	"github.com/google/uuid"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/grpc"
)

//StatusArgs arguments for status cluster
type StatusArgs struct {
	Address string
	Port    int32
}

//ExecuteCmdStatus helper function to execute grpc call to get the status
func ExecuteCmdStatus(args *StatusArgs) (*pb.ClusterStatusReply, error) {
	peerConn, rpcClient, err := rpc.GetPeerConn(args.Address, args.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to get peer connection, error = %w", err)
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	clusterStatus, err := rpcClient.ClusterStatus(context.Background(), &pb.ClusterStatusRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ClusterStatus, error = %w", err)
	}
	return clusterStatus, err
}

//CreateCommandArgs argument structure for this command
type CreateCommandArgs struct {
	Address string `json:"address,omitempty"`
	Port    int32  `json:"port,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

//ExecuteCmdCreateStore helper function to executed create store command
func ExecuteCmdCreateStore(args *CreateCommandArgs) error {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return err
	}

	peerConn, rpcClient, err1 := rpc.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return err1
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := &pb.StoreCreateRequest{
		Name: args.Name,
		Id:   uuid.New().String(),
		Type: pb.StoreType_Memory,
	}
	if args.Type == "sqlite" {
		req.Type = pb.StoreType_Sqlite
	}

	_, err = rpcClient.CreateStore(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

//GetArgs argument structure for this command
type GetArgs struct {
	Address   string
	Port      int32
	Key       string
	StoreName string
}

//ExecuteCmdGet helper function to perform put command
func ExecuteCmdGet(args *GetArgs) (*pb.GetReply, error) {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return nil, err
	}
	peerConn, rpcClient, err := rpc.GetPeerConn(leader.Address, leader.Port)
	if err != nil {
		return nil, err
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := &pb.GetRequest{
		StoreName: args.StoreName,
		Key:       args.Key,
	}
	reply, err := rpcClient.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

//ListArgs argument structure for this command
type ListArgs struct {
	Address string
	Port    int32
}

//ExecuteCmdListStore executes the list command for a store
func ExecuteCmdListStore(args *ListArgs) (*pb.ListStoreReply, error) {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return nil, err
	}

	peerConn, rpcClient, err := rpc.GetPeerConn(leader.Address, leader.Port)
	if err != nil {
		return nil, err
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	storeList, err := rpcClient.ListStores(context.Background(), &pb.ListStoreRequest{})
	if err != nil {
		return nil, err
	}
	return storeList, err
}

//PutArgs argument structure for this command
type PutArgs struct {
	Address   string `json:"address,omitempty"`
	Port      int32  `json:"port,omitempty"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	StoreName string `json:"storeName"`
}

//ExecuteCmdPut helper function to perform put command
func ExecuteCmdPut(args *PutArgs) error {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return err
	}

	log.Printf("sending request to leader [%s:%s:%d]", leader.Name, leader.Address, leader.Port)

	peerConn, rpcClient, err1 := rpc.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return err1
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := &pb.PutRequest{
		StoreName: args.StoreName,
		Key:       args.Key,
		Value:     args.Value,
		Id:        uuid.New().String(),
	}

	_, err = rpcClient.Put(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

//QueryArgs arguments for the query command
type QueryArgs struct {
	Address   string `json:"address,omitempty"`
	Port      int32  `json:"port,omitempty"`
	SQLStr    string `json:"sqlite"`
	StoreName string `json:"storeName"`
}

//ExecuteCmdQuery executes the query on the store
func ExecuteCmdQuery(args *QueryArgs) (*pb.ExecuteQueryReply, error) {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return nil, err
	}

	peerConn, rpcClient, err1 := rpc.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return nil, err1
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := pb.ExecuteQueryRequest{
		Id:        uuid.New().String(),
		StoreName: args.StoreName,
		SqlQuery:  args.SQLStr,
	}

	resp, err1 := rpcClient.ExecuteQuery(context.Background(), &req)
	if err1 != nil {
		return nil, err1
	}
	return resp, nil
}

//ExecuteCmdSQLQuery executes the query on the store
func ExecuteCmdSQLQuery(args *QueryArgs) (*pb.QueryReply, error) {
	leader, err := GetLeader(args.Address, args.Port)
	if err != nil {
		return nil, err
	}

	peerConn, rpcClient, err1 := rpc.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return nil, err1
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := pb.QueryRequest{
		Id:        uuid.New().String(),
		StoreName: args.StoreName,
		SqlQuery:  args.SQLStr,
	}

	resp, err1 := rpcClient.Query(context.Background(), &req)
	if err1 != nil {
		return nil, err1
	}
	return resp, nil
}
