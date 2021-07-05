package main

import (
	"github.com/davinash/yados/commands/yados"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yados",
	Short: "YADOS CLI",
	Long: `
YADOS cli to interact with st cluster`,
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
	yados.AddCommands(rootCmd)
	Execute()
}
