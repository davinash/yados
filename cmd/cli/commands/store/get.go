package store

import (
	"context"
	"fmt"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//GetArgs argument structure for this command
type GetArgs struct {
	Address   string
	Port      int32
	Key       string
	StoreName string
}

//ExecuteCmdGet helper function to perform put command
func ExecuteCmdGet(args *GetArgs) (*pb.GetReply, error) {
	leader, err := server.GetLeader(args.Address, args.Port)
	if err != nil {
		return nil, err
	}
	peerConn, rpcClient, err := server.GetPeerConn(leader.Address, leader.Port)
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

//CreateGetCommand cobra command for listing stores
func CreateGetCommand(rootCmd *cobra.Command) {
	getArg := GetArgs{}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a value for a key from a store ",
		RunE: func(cmd *cobra.Command, args []string) error {
			reply, err := ExecuteCmdGet(&getArg)
			if err != nil {
				return err
			}
			fmt.Println(reply.Value)
			return nil
		},
	}
	cmd.Flags().StringVar(&getArg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&getArg.Port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&getArg.Key, "key", "", "Key name")
	err := cmd.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&getArg.StoreName, "store-name", "", "store in which get operations to be "+
		"performed")
	err = cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
