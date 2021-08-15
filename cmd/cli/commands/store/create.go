package store

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//CreateCommandArgs argument structure for this command
type CreateCommandArgs struct {
	Address string
	Port    int32
	Name    string
}

//CreateCommandExecute helper function to executed create store command
func CreateCommandExecute(args *CreateCommandArgs) error {
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

	req := &pb.StoreCreateRequest{Name: args.Name}
	marshal, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	_, err = rpcClient.RunCommand(context.Background(), &pb.CommandRequest{
		Id:      uuid.New().String(),
		Args:    marshal,
		CmdType: pb.CommandType_CreateStore,
	})
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
			return CreateCommandExecute(&arg)
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
