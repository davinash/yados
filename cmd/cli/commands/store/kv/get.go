package kv

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateGetCommand cobra command for listing stores
func CreateGetCommand(rootCmd *cobra.Command) {
	var address string
	var port int32
	getArg := server.GetArgs{}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a value for a key from a store ",
		Long: `
### Get value of key from KV Store
yadosctl store kv get --store-name KvStore --key key1
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reply, err := server.ExecuteCmdGet(&getArg, address, port)
			if err != nil {
				return err
			}
			fmt.Println(reply.Value)
			return nil
		},
	}
	cmd.Flags().StringVar(&address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

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
