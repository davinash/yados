package controller

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/davinash/yados/internal/controller"
	"github.com/spf13/cobra"
)

//StartCommands command for starting a server
func StartCommands(rootCmd *cobra.Command) {
	var address string
	var port int32
	cmd := &cobra.Command{
		Use:   "start",
		Short: "create and start a new controller",
		Long: `
For Example:
### Start controller on default port ( 127.0.0.1:9090 )
yadosctl controller start

### Start controller on specified address and port
yadosctl controller start --listen-address 127.0.0.1 --port 9090
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			oSSignalCh := make(chan os.Signal, 1)
			signal.Notify(oSSignalCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			c := controller.NewController(address, port, "info")
			c.Start()

			<-oSSignalCh
			return nil
		},
	}
	cmd.Flags().StringVar(&address, "listen-address", "127.0.0.1",
		"Listen Address on which controller will listen\n"+
			"Can usually be left blank. Otherwise, use IP address or host Name \n"+
			"that other server nodes use to connect")

	cmd.Flags().Int32Var(&port, "port", 9090, "Port to use for communication")
	rootCmd.AddCommand(cmd)
}
