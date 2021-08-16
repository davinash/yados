package store

import (
	"context"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

//PutArgs argument structure for this command
type PutArgs struct {
	Address   string
	Port      int32
	Key       string
	Value     string
	StoreName string
}

//ExecutePutCommand helper function to perform put command
func ExecutePutCommand(args *PutArgs) error {
	leader, err := server.GetLeader(args.Address, args.Port)
	if err != nil {
		return err
	}

	peerConn, rpcClient, err := server.GetPeerConn(leader.Address, leader.Port)
	if err != nil {
		return err
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
	}
	marshal, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	_, err = rpcClient.RunCommand(context.Background(), &pb.CommandRequest{
		Id:      uuid.New().String(),
		Args:    marshal,
		CmdType: pb.CommandType_Put,
	})
	if err != nil {
		return err
	}
	return nil
}

//CreatePutCommand cobra command for listing stores
func CreatePutCommand(rootCmd *cobra.Command) {
	putArg := PutArgs{}
	cmd := &cobra.Command{
		Use:   "put",
		Short: "put a key/value in a store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecutePutCommand(&putArg)
		},
	}
	cmd.Flags().StringVar(&putArg.Address, "Address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&putArg.Port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&putArg.Key, "key", "", "Key name")
	err := cmd.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&putArg.Value, "value", "", "value for a key")
	err = cmd.MarkFlagRequired("value")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&putArg.StoreName, "store-name", "", "store in which put operations to be "+
		"performed")
	err = cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
