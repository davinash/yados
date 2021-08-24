package server

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//StartSrvArgs represents the Start server Arguments
type StartSrvArgs struct {
	name     string
	address  string
	port     int32
	peers    []string
	logLevel string
	logDir   string
	httpPort int
}

//StartCommands command for starting a server
func StartCommands(rootCmd *cobra.Command) {
	srvStartArgs := StartSrvArgs{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "create and start a new server",
		RunE: func(cmd *cobra.Command, args []string) error {
			oSSignalCh := make(chan os.Signal, 1)
			signal.Notify(oSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			srvArgs := &server.NewServerArgs{
				Name:     srvStartArgs.name,
				Address:  srvStartArgs.address,
				Port:     srvStartArgs.port,
				Loglevel: srvStartArgs.logLevel,
				LogDir:   srvStartArgs.logDir,
				HTTPPort: srvStartArgs.httpPort,
			}
			srv, err := server.NewServer(srvArgs)
			if err != nil {
				return err
			}
			peers := make([]*pb.Peer, 0)
			for _, p := range srvStartArgs.peers {
				split := strings.Split(p, ":")
				if len(split) != 3 {
					return fmt.Errorf("invalid format for peers, use <name:ip-address:port>")
				}
				port, err := strconv.Atoi(split[2])
				if err != nil {
					return fmt.Errorf("invalid format for peers, use <name:ip-address:port>")
				}
				peer := &pb.Peer{
					Name:    split[0],
					Address: split[1],
					Port:    int32(port),
				}
				peers = append(peers, peer)
			}
			err = srv.Serve(peers)
			if err != nil {
				return err
			}
			<-oSSignalCh
			err = srv.Stop()
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&srvStartArgs.name, "name", "", "Name of the server, Name must be unique in a cluster")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&srvStartArgs.address, "listen-address", "127.0.0.1",
		"Listen Address on which server will listen\n"+
			"Can usually be left blank. Otherwise, use IP address or host Name \n"+
			"that other server nodes use to connect to the new server")

	cmd.Flags().Int32Var(&srvStartArgs.port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringSliceVar(&srvStartArgs.peers, "peer", []string{},
		"peer to join <name:ip-address:port>, "+
			"use multiple of this flag if want to join with multiple peers")

	cmd.Flags().StringVar(&srvStartArgs.logLevel, "log-level", "info", "Log level "+
		"[info|debug|warn|trace|error]")

	cmd.Flags().StringVar(&srvStartArgs.logDir, "log-dir", "info",
		"Location for replicated log storage")
	err = cmd.MarkFlagRequired("log-dir")
	if err != nil {
		panic(err)
	}
	cmd.Flags().IntVar(&srvStartArgs.httpPort, "http-port", -1, "Port to use for http server")

	rootCmd.AddCommand(cmd)
}
