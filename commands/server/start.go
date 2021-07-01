package server

import (
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

func AddNodeStartCmd(parentCmd *cobra.Command) {
	var serverName string
	var address string
	var port int
	var clusterName string
	var withPeer bool
	var peerAddress string
	var peerPort int
	cmd := &cobra.Command{
		Use:   "start",
		Short: "startHttpServer a server",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.CreateNewServer(serverName, address, port, clusterName,
				withPeer, peerAddress, peerPort)
			if err != nil {
				return err
			}
			return srv.StartAndWait()
		},
	}
	cmd.Flags().StringVar(&serverName, "name", "", "Name of the server")
	_ = cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(&address, "listen-address", "127.0.0.1", "Listen Address on which server will listen\n"+
		"Can usually be left blank. Otherwise, use IP address or host name \nthat other server nodes use to connect to the new server")

	cmd.Flags().IntVar(&port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of the cluster\n"+
		"Same cluster name need to use for other servers \nin order to participate in the cluster")
	_ = cmd.MarkFlagRequired("clusterName")

	//--with-peer --peer-address --peer-port
	cmd.Flags().BoolVar(&withPeer, "with-peer", false, "Use this flag if server need to part of cluster")
	cmd.Flags().StringVar(&peerAddress, "peer-address", "", "IP address or host name of the peer")
	cmd.Flags().IntVar(&peerPort, "peer-port", -1, "Port to use for communication with peer")

	parentCmd.AddCommand(cmd)
}
