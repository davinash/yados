package server

import (
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

// AddServerStartCmd Cobra command implementation for Starting a node for the cluster
func AddServerStartCmd(parentCmd *cobra.Command) {
	var serverName string
	var address string
	var port int32
	var peers []string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start YADOS YadosServer",
		Long: "" +
			"Example : yados server start --name server1\n" +
			"          yados server start --name server1 --port 9191 --listen-address 127.0.0.1",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.CreateNewServer(serverName, address, port)
			if err != nil {
				return err
			}
			return srv.StartAndWait(peers)
		},
	}
	cmd.Flags().StringVar(&serverName, "name", "", "Name of the server")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(&address, "listen-address", "127.0.0.1", "Listen Address on which server will listen\n"+
		"Can usually be left blank. Otherwise, use IP address or host name \nthat other server nodes use to connect to the new server")

	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringSliceVar(&peers, "peer", []string{}, "peer to join <ip-address:port>, "+
		"use multiple of this flag if want to join with multiple peers")

	parentCmd.AddCommand(cmd)
}
