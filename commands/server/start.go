package server

import (
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AddServerStartCmd Cobra command implementation for Starting a node for the cluster
func AddServerStartCmd(parentCmd *cobra.Command) {
	var serverName string
	var address string
	var port int32
	var withPeer bool
	var peerAddress string
	var peerPort int32
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start YADOS YadosServer",
		Long: "" +
			"Example : yados server start --name server1\n" +
			"          yados server start --name server1 --port 9191 --listen-address 127.0.0.1",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			isWithPeer := false
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if f.Name == "with-peer" {
					isWithPeer = true
				}
			})
			if isWithPeer {
				isPeerAddressPassed := false
				cmd.Flags().Visit(func(f *pflag.Flag) {
					if f.Name == "peer-address" {
						isPeerAddressPassed = true
					}
				})
				if !isPeerAddressPassed {
					return fmt.Errorf("peer-address is missing, check help")
				}
			}

			if isWithPeer {
				isPeerPort := false
				cmd.Flags().Visit(func(f *pflag.Flag) {
					if f.Name == "peer-port" {
						isPeerPort = true
					}
				})
				if !isPeerPort {
					return fmt.Errorf("peer-port is missing, check help")
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.CreateNewServer(serverName, address, port)
			if err != nil {
				return err
			}
			return srv.StartAndWait(withPeer, peerAddress, peerPort)
		},
	}
	cmd.Flags().StringVar(&serverName, "name", "", "Name of the server")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(&address, "listen-address", "127.0.0.1", "Listen Address on which server will listen\n"+
		"Can usually be left blank. Otherwise, use IP address or host name \nthat other server nodes use to connect to the new server")

	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

	cmd.Flags().BoolVar(&withPeer, "with-peer", false, "Use this flag if server need to part of cluster")
	cmd.Flags().StringVar(&peerAddress, "peer-address", "", "IP address or host name of the peer")
	cmd.Flags().Int32Var(&peerPort, "peer-port", -1, "Port to use for communication with peer")

	parentCmd.AddCommand(cmd)
}
