package main

import (
	"github.com/davinash/yados/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yados",
	Short: "YADOS CLI",
	Long: `
YADOS cli to interact with YADOS cluster`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func initConfig() {
}

func init() {
	cobra.OnInitialize(initConfig)
}

func main() {
	commands.AddCommands(rootCmd)
	Execute()
}
