package store

import (
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateCommand WAL create command
func CreateCommand(rootCmd *cobra.Command) {
	arg := server.CreateCommandArgs{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new store",
		Long: `
For Example:

Creating store with default type
./yadosctl store create --store-name store1

Creating store with specifying type 
./yadosctl store create --store-name store1 --store-type kv

Creating store with Sqlite. Internally yados will be using sqlite to store the data
./yadosctl store create --store-name SqlStore1 --store-type sqlite
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.ExecuteCmdCreateStore(&arg)
		},
	}
	cmd.Flags().StringVar(&arg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&arg.Port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&arg.Name, "store-name", "", "Name of the store to create")
	err := cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&arg.Type, "store-type", "memory", "type of the store to create [ memory | sqlite ]")

	rootCmd.AddCommand(cmd)
}
