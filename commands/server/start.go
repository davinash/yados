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
	var hcTimeDuration int
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start YADOS YadosServer",
		Long: "" +
			"Example : \n\n" +
			"Start server with default options\n" +
			"    yados server start --name server1\n" +
			"Start server with specific IP address and port\n" +
			"    yados server start --name server1 --port 9191 --listen-address 127.0.0.1\n" +
			"Start server with more than one peer to form a cluster\n" +
			"    yados server start --port 9193 --peer 127.0.0.1:9191 --name Server3 --peer 127.0.0.1:9192\n\n",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.CreateNewServer(serverName, address, port)
			if err != nil {
				return err
			}
			srv.SetHealthCheckDuration(hcTimeDuration)
			err = srv.Start(peers)
			if err != nil {
				return err
			}
			<-srv.OSSignalCh
			return nil
		},
	}
	cmd.Flags().StringVar(&serverName, "name", "", "Name of the server, name must be unique in a cluster")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(&address, "listen-address", "127.0.0.1", "Listen Address on which server will listen\n"+
		"Can usually be left blank. Otherwise, use IP address or host name \nthat other server nodes use to connect to the new server")

	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringSliceVar(&peers, "peer", []string{}, "peer to join <ip-address:port>, "+
		"use multiple of this flag if want to join with multiple peers")

	cmd.Flags().IntVar(&hcTimeDuration, "health-check-duration", 30, "Health check duration in seconds")

	parentCmd.AddCommand(cmd)
}
