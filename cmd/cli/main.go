package main

import (
	"github.com/davinash/yados/commands/cli"
	"github.com/spf13/cobra"
)

var (
	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "yadosctl",
	Short: "yadosctl",
	Long: `
yadosct cli to interact with YADOS cluster`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func initConfig() {
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
}

func main() {
	cli.AddCommands(rootCmd)
	Execute()
}
