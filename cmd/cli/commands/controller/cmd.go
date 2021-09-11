package controller

import (
	"github.com/spf13/cobra"
)

//NewControllerCommands creates cobra command for controller
func NewControllerCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "YADOS controller manage and monitor commands",
	}
	StartCommands(cmd)
	return cmd
}
