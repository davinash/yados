package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//ListArgs argument structure for this command
type ListArgs struct {
	Address string
	Port    int32
}

//ExecuteCmdListStore executes the list command for a store
func ExecuteCmdListStore(args *ListArgs) (*pb.ListStoreReply, error) {
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

	storeList, err := rpcClient.ListStores(context.Background(), &pb.ListStoreRequest{})
	if err != nil {
		return nil, err
	}
	return storeList, err
}

//CreateListCommand cobra command for listing stores
func CreateListCommand(rootCmd *cobra.Command) {
	listArg := ListArgs{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stores",
		RunE: func(cmd *cobra.Command, args []string) error {

			storeList, err := ExecuteCmdListStore(&listArg)
			if err != nil {
				return err
			}
			marshal, err := json.MarshalIndent(storeList.Name, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(marshal))
			return nil
		},
	}
	cmd.Flags().StringVar(&listArg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&listArg.Port, "port", 9191, "Port to use for communication")

	rootCmd.AddCommand(cmd)

}
