package store

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// AddCreateStoreCommand Cobra command for creating a store
func AddCreateStoreCommand(parentCmd *cobra.Command) {
	var storeName string
	var replication int32
	var serverAddress string
	var port int32
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create store",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, p, err := server.GetPeerConn(serverAddress, port)
			if err != nil {
				return err
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
				}
			}(conn)

			_, err = p.CreateStore(context.Background(), &pb.StoreCreateRequest{
				Name:        storeName,
				Replication: replication,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&storeName, "name", "", "Name of the store, should be unique in the cluster")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().Int32Var(&replication, "replication", 1, "Replication factor for the store")
	cmd.Flags().StringVar(&serverAddress, "server", "127.0.0.1", "IP address or host name")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	parentCmd.AddCommand(cmd)
}
