package store

import (
	"encoding/json"
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateListCommand cobra command for listing stores
func CreateListCommand(rootCmd *cobra.Command) {
	listArg := server.ListArgs{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stores",
		RunE: func(cmd *cobra.Command, args []string) error {

			storeList, err := server.ExecuteCmdListStore(&listArg)
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
