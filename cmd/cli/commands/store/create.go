package store

import (
	"context"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//CreateCommandArgs argument structure for this command
type CreateCommandArgs struct {
	address string
	port    int32
	name    string
}

//CreateCommand PLog create command
func CreateCommand(rootCmd *cobra.Command) {
	arg := CreateCommandArgs{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new store",
		RunE: func(cmd *cobra.Command, args []string) error {
			peerConn, rpcClient, err := server.GetPeerConn(arg.address, arg.port)
			if err != nil {
				return err
			}
			defer func(peerConn *grpc.ClientConn) {
				err := peerConn.Close()
				if err != nil {
					log.Printf("failed to close the connection, error = %v\n", err)
				}
			}(peerConn)

			arg := &pb.StoreCreateRequest{Name: arg.name}
			_, err = rpcClient.CreateStore(context.Background(), arg)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&arg.address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&arg.port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&arg.name, "name", "", "Name of the store to create")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
