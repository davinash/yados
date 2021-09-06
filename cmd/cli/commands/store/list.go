package store

import (
	"encoding/json"
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateListCommand cobra command for listing stores
func CreateListCommand(rootCmd *cobra.Command) {
	var address string
	var port int32
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stores",
		Long: `

### To list all the stores
yadosctl store list
`,
		RunE: func(cmd *cobra.Command, args []string) error {

			storeList, err := server.ExecuteCmdListStore(address, port)
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
	cmd.Flags().StringVar(&address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	rootCmd.AddCommand(cmd)

}
