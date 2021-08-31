package sql

import (
	"encoding/json"
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateExecuteQueryCommand cobra command construction of query
func CreateExecuteQueryCommand(rootCmd *cobra.Command) {
	queryArg := server.QueryArgs{}
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "execute sql query on the store ( DDL ) ",
		RunE: func(cmd *cobra.Command, args []string) error {
			reply, err := server.ExecuteQuery(&queryArg)
			if err != nil {
				return err
			}
			replyBytes, err := json.Marshal(reply)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", string(replyBytes))
			return nil
		},
	}
	cmd.Flags().StringVar(&queryArg.Address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&queryArg.Port, "port", 9191, "Port to use for communication")

	cmd.Flags().StringVar(&queryArg.SQLStr, "sql", "", "SQL statement to execute")
	err := cmd.MarkFlagRequired("sql")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVar(&queryArg.StoreName, "store-name", "", "Name of the store to create")
	err = cmd.MarkFlagRequired("store-name")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(cmd)

}
