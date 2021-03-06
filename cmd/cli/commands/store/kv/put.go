package kv

import (
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreatePutCommand cobra command for listing stores
func CreatePutCommand(rootCmd *cobra.Command) {
	var address string
	var port int32
	putArg := server.PutArgs{}
	cmd := &cobra.Command{
		Use:   "put",
		Short: "puts a key/value in a store",
		Long: `
### Insert a value in a kv store
yadosctl store kv put --store-name KvStore --key key1 --value value1
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ExecuteCmdPut(&putArg, address, port)
		},
	}
	cmd.Flags().StringVar(&address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

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
