package sql

import (
	"encoding/json"
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//ColumnV column value structure
type ColumnV struct {
	Name  string
	Value string
}

//TRow represents the row
type TRow struct {
	Row []*ColumnV
}

//CreateQueryCommand cobra command construction of query
func CreateQueryCommand(rootCmd *cobra.Command) {
	queryArg := server.QueryArgs{}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "execute sql query on the store ( DML )",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := server.ExecuteCmdSQLQuery(&queryArg)
			if err != nil {
				//return err
				panic(err)
			}
			rows := make([]*TRow, 0)

			for _, r := range result.Rows {
				tRow := TRow{
					Row: make([]*ColumnV, 0),
				}
				for _, rr := range r.Row {
					var value interface{}
					err := json.Unmarshal(rr.Value, &value)
					if err != nil {
						//return err
						panic(err)
					}
					tRow.Row = append(tRow.Row, &ColumnV{
						Name:  rr.Name,
						Value: fmt.Sprintf("%v", value),
					})
				}
				rows = append(rows, &tRow)
			}
			bytes, err := json.MarshalIndent(rows, "", "  ")
			if err != nil {
				//return err
				panic(err)
			}
			fmt.Println(string(bytes))
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
