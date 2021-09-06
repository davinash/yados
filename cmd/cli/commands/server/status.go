package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//StatusCommands command for showing the status of one more servers
func StatusCommands(rootCmd *cobra.Command) {
	statusArgs := &server.StatusArgs{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status of the servers in cluster",
		Long: `
### Get the status of the cluster
yadosctl server status
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterStatus, err := server.ExecuteCmdStatus(statusArgs.Address, statusArgs.Port)
			if err != nil {
				return err
			}
			marshal, err := json.MarshalIndent(clusterStatus.PeerStatus, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(marshal))
			os.Exit(0)
			return nil
		},
	}
	cmd.Flags().StringVar(&statusArgs.Address, "address", "127.0.0.1", "server to connect in the cluster")
	cmd.Flags().Int32Var(&statusArgs.Port, "port", 9191, "Port to use for communication")
	rootCmd.AddCommand(cmd)
}
