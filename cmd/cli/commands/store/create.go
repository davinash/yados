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

//CreateCommandArgs argument structure for this command
type CreateCommandArgs struct {
	Address string
	Port    int32
	Name    string
}

//ExecuteCmdCreateStore helper function to executed create store command
func ExecuteCmdCreateStore(args *CreateCommandArgs) error {
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

	req := &pb.StoreCreateRequest{
		Name: args.Name,
		Id:   uuid.New().String(),
	}

	_, err = rpcClient.CreateStore(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

//CreateCommand PLog create command
func CreateCommand(rootCmd *cobra.Command) {
	arg := CreateCommandArgs{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteCmdCreateStore(&arg)
		},
	}
	cmd.Flags().StringVar(&arg.Address, "Address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&arg.Port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&arg.Name, "name", "", "Name of the store to create")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
