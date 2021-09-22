package server

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/status"

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
func ExecuteCmdStatus(address string, port int32) (*pb.ClusterStatusReply, error) {
	peerConn, rpcClient, err := rpc.GetPeerConn(address, port)
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
		e, _ := status.FromError(err)
		e.Message()
		return &pb.ClusterStatusReply{}, fmt.Errorf("ClusterStatus failed, Error= %s", e.Message())
	}
	return clusterStatus, err
}

//CreateCommandArgs argument structure for this command
type CreateCommandArgs struct {
	Name string `json:"name"`
}

//ExecuteCmdCreateStore helper function to executed create store command
func ExecuteCmdCreateStore(args *CreateCommandArgs, address string, port int32) (*pb.StoreCreateReply, error) {
	reply := &pb.StoreCreateReply{}

	leader, err := GetLeader(address, port)
	if err != nil {
		return reply, err
	}

	peerConn, rpcClient, err := rpc.GetPeerConn(leader.Address, leader.Port)
	if err != nil {
		return reply, err
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
	}

	resp, err := rpcClient.CreateStore(context.Background(), req)
	if err != nil {
		e, _ := status.FromError(err)
		e.Message()
		return reply, fmt.Errorf("CreateStore failed, Error= %s", e.Message())
	}

	return resp, nil
}

//ExecuteCmdListStore executes the list command for a store
func ExecuteCmdListStore(address string, port int32) (*pb.ListStoreReply, error) {
	leader, err := GetLeader(address, port)
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
		e, _ := status.FromError(err)
		e.Message()
		return &pb.ListStoreReply{}, fmt.Errorf("ListStores failed, Error= %s", e.Message())
	}
	return storeList, err
}

//QueryArgs arguments for the query command
type QueryArgs struct {
	SQLStr    string `json:"sql"`
	StoreName string `json:"storeName"`
}

//ExecuteCmdQuery executes the query on the store
func ExecuteCmdQuery(args *QueryArgs, address string, port int32) (*pb.ExecuteQueryReply, error) {
	leader, err := GetLeader(address, port)
	if err != nil {
		return nil, err
	}

	peerConn, rpcClient, err1 := rpc.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return &pb.ExecuteQueryReply{}, err1
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

	resp, err := rpcClient.ExecuteQuery(context.Background(), &req)
	if err != nil {
		e, _ := status.FromError(err)
		e.Message()
		return &pb.ExecuteQueryReply{}, fmt.Errorf("ExecuteQuery failed, Error= %s", e.Message())
	}
	return resp, err
}

//ExecuteCmdSQLQuery executes the query on the store
func ExecuteCmdSQLQuery(args *QueryArgs, address string, port int32) (*pb.QueryReply, error) {
	leader, err := GetLeader(address, port)
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

	resp, err := rpcClient.Query(context.Background(), &req)
	if err != nil {
		e, _ := status.FromError(err)
		e.Message()
		return &pb.QueryReply{}, fmt.Errorf("query failed, Error= %s", e.Message())
	}
	return resp, nil
}
