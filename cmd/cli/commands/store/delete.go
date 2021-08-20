package store

import (
	"context"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//DeleteArgs arguments for delete
type DeleteArgs struct {
	Address   string
	Port      int32
	StoreName string
}

//ExecuteCmdDeleteStore execute the delete command
func ExecuteCmdDeleteStore(args *DeleteArgs) error {
	leader, err := server.GetLeader(args.Address, args.Port)
	if err != nil {
		return err
	}
	peerConn, rpcClient, err1 := server.GetPeerConn(leader.Address, leader.Port)
	if err1 != nil {
		return err1
	}
	defer func(peerConn *grpc.ClientConn) {
		err := peerConn.Close()
		if err != nil {
			log.Printf("failed to close the connection, error = %v\n", err)
		}
	}(peerConn)

	req := &pb.StoreDeleteRequest{
		StoreName: args.StoreName,
		Id:        uuid.New().String(),
	}

	_, err = rpcClient.DeleteStore(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

//CreateDeleteCommand cobra command construction for delete operation
func CreateDeleteCommand(rootCmd *cobra.Command) {
	deleteArg := DeleteArgs{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete stores",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ExecuteCmdDeleteStore(&deleteArg)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&deleteArg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&deleteArg.Port, "port", 9191, "Port to use for communication")
	cmd.Flags().StringVar(&deleteArg.StoreName, "store-name", "", "Name of the store to delete")
	err := cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
