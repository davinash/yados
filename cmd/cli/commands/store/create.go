package store

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateCommand WAL create command
func CreateCommand(rootCmd *cobra.Command) {
	var address string
	var port int32
	arg := server.CreateCommandArgs{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create new store",
		Long: `

### Create Store of type sqlite
yadosctl store create --store-name store-sqlite --store-type sqlite

### Create store of type KV
yadosctl store create --store-name store-kv
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := server.ExecuteCmdCreateStore(&arg, address, port)
			if err != nil {
				return err
			}
			fmt.Println(resp.Msg)
			return nil
		},
	}
	cmd.Flags().StringVar(&address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&arg.Name, "store-name", "", "Name of the store to create")
	err := cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)
}
