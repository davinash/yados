package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

//ListArg argument structure for list
type ListArg struct {
	address string
	port    int32
}

//ListCommands cobra command for listing the servers
func ListCommands(rootCmd *cobra.Command) {
	arg := &ListArg{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the servers in the cluster",
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

			reply, err := rpcClient.ListMembers(context.Background(), &pb.ListMembersRequest{})
			if err != nil {
				return err
			}
			marshal, err := json.MarshalIndent(reply.Peers, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(marshal))
			os.Exit(0)
			return nil
		},
	}
	cmd.Flags().StringVar(&arg.address, "address", "127.0.0.1", "server to connect in the cluster")
	cmd.Flags().Int32Var(&arg.port, "port", 9191, "Port to use for communication")

	rootCmd.AddCommand(cmd)
}
