package main

import (
	"github.com/davinash/yados/cmd/cli/commands/server"
	"github.com/davinash/yados/cmd/cli/commands/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yadosctl",
	Short: "yadosctl",
	Long: `
yadosct is a cli to manage and monitor Yet Another Distributed Object PLog (YADOS) cluster`,
}

//Execute main driver function form Cobra commands
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func initConfig() {
}

func init() {
	cobra.OnInitialize(initConfig)
}

func main() {
	server.AddCommands(rootCmd)
	store.AddCommands(rootCmd)
	Execute()
}
