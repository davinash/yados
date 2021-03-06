package sql

import (
	"encoding/json"
	"fmt"

	"github.com/davinash/yados/internal/server"
	"github.com/spf13/cobra"
)

//CreateExecuteQueryCommand cobra command construction of query
func CreateExecuteQueryCommand(rootCmd *cobra.Command) {
	var address string
	var port int32
	queryArg := server.QueryArgs{}
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "execute sql query on the store",
		Long: `
### Execute SQL command
#### Create a table
yadosctl store sqlite execute --store-name SqlStore1 --sql "create table employee(empid integer,name varchar(20),title varchar(10))"
#### Insert
yadosctl store sqlite execute --store-name SqlStore1 --sql "insert into employee values(101,'John Smith','CEO')" 
#### Insert
yadosctl store sqlite execute --store-name SqlStore1 --sql "insert into employee values(102,'Raj Reddy','Sysadmin')" 
#### Insert
yadosctl store sqlite execute --store-name SqlStore1 --sql "insert into employee values(103,'Jason Bourne','Developer')" 
#### Insert
yadosctl store sqlite execute --store-name SqlStore1 --sql "insert into employee values(104,'Jane Smith','Sale Manager')"
#### Insert
yadosctl store sqlite execute --store-name SqlStore1 --sql "insert into employee values(105,'Rita Patel','DBA')"
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			reply, err := server.ExecuteCmdQuery(&queryArg, address, port)
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
	cmd.Flags().StringVar(&address, "address", "127.0.0.1", "Server to connect in the cluster")
	cmd.Flags().Int32Var(&port, "port", 9191, "Port to use for communication")

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
