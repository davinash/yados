package server

import (
	"github.com/spf13/cobra"
)

//NewServerCommands creates cobra command for server
func NewServerCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "YADOS server manage and monitor commands",
	}
	StartCommands(cmd)
	StatusCommands(cmd)
	return cmd
}
