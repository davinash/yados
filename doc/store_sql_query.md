## store sql query

execute sql query on the store

### Synopsis


### Query using SQL
yadosctl store sqlite query --store-name SqlStore1 --sql "select * from employee" --address 192.28.0.2


```
store sql query [flags]
```

### Options

```
      --address string      Server to connect in the cluster (default "127.0.0.1")
  -h, --help                help for query
      --port int32          Port to use for communication (default 9191)
      --sql string          SQL statement to execute
      --store-name string   Name of the store to create
```

### SEE ALSO

* [store sql](store_sql.md)	 - manage and monitor store for sql store ( sqlite )

###### Auto generated by spf13/cobra on 7-Sep-2021
