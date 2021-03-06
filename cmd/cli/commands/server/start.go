package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//StartSrvArgs represents the Start server Arguments
type StartSrvArgs struct {
	name    string
	address string
	port    int32
	//peers       []string
	logLevel    string
	walDir      string
	httpPort    int
	controllers []string
}

//StartCommands command for starting a server
func StartCommands(rootCmd *cobra.Command) {
	srvStartArgs := StartSrvArgs{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "create and start a new server",
		Long: `Example of commands to start the server and create the cluster

### Starting server with default options
yadosctl server start --name Server1 --wal-dir /tmp --controller 127.0.0.1:9090

### Starting server with options
yadosctl server start --name Server1 --listen-address 127.0.0.1 --wal-dir /tmp --port 9191 --log-level info

### Starting server with options with http server
yadosctl server start --name Server1 --listen-address 127.0.0.1 --wal-dir /tmp --port 9191 --log-level info --http-port 8181

### Starting second server and join the cluster
yadosctl server start --name server2 --listen-address 127.0.0.1 --wal-dir /tmp --port 9192 --controller 127.0.0.1:9090

### Starting third server and join the cluster
yadosctl server start --name server3 --listen-address 127.0.0.1 --wal-dir /tmp --port 9193 --controller 127.0.0.1:9090
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			oSSignalCh := make(chan os.Signal, 1)
			signal.Notify(oSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			srvArgs := &server.NewServerArgs{
				Name:        srvStartArgs.name,
				Address:     srvStartArgs.address,
				Port:        srvStartArgs.port,
				Loglevel:    srvStartArgs.logLevel,
				WalDir:      srvStartArgs.walDir,
				HTTPPort:    srvStartArgs.httpPort,
				Controllers: srvStartArgs.controllers,
			}
			srv, err := server.NewServer(srvArgs)
			if err != nil {
				return err
			}
			err = srv.Serve()
			if err != nil {
				return err
			}
			<-oSSignalCh
			err = srv.Stop()
			if err != nil {
				return fmt.Errorf("failed to stop the server, error = %w", err)
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

	cmd.Flags().StringVar(&srvStartArgs.logLevel, "log-level", "info", "Log level "+
		"[info|debug|warn|trace|error]")

	cmd.Flags().StringVar(&srvStartArgs.walDir, "wal-dir", "info",
		"Location for replicated write ahead log")
	err = cmd.MarkFlagRequired("wal-dir")
	if err != nil {
		panic(err)
	}
	cmd.Flags().IntVar(&srvStartArgs.httpPort, "http-port", -1, "Port to use for http server")

	cmd.Flags().StringSliceVar(&srvStartArgs.controllers, "controller", []string{},
		"controller <ip-address:port>, "+
			"use multiple of this flag if want to join with multiple controller")

	rootCmd.AddCommand(cmd)
}
