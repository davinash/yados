package store

import (
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreatePutCommand cobra command for listing stores
func CreatePutCommand(rootCmd *cobra.Command) {
	putArg := server.PutArgs{}
	cmd := &cobra.Command{
		Use:   "put",
		Short: "put a key/value in a store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ExecuteCmdPut(&putArg)
		},
	}
	cmd.Flags().StringVar(&putArg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&putArg.Port, "port", 9191, "Port to use for communication")

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
