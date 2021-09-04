# yados
<img src="https://github.com/davinash/yados/blob/main/logo.jpg" align="center"/>

Yet Another Distributed Object Store

[![Go Actions Status](https://github.com/davinash/yados/workflows/Go/badge.svg)](https://github.com/davinash/yados/actions)

## Starting a server and joining the cluster
```shell
Starting server with default options
./yadosctl server start --name Server1 --wal-dir /tmp

Starting server with options
./yadosctl server start --name Server1 --listen-address 127.0.0.1 --wal-dir /tmp --port 9191 --log-level info

Starting server with options with http server
./yadosctl server start --name Server1 --listen-address 127.0.0.1 --wal-dir /tmp --port 9191 --log-level info --http-port 8181

Starting second server and join the cluster
./yadosctl server start --name server2 --listen-address 127.0.0.1 --wal-dir /tmp --port 9192  --peer server1:127.0.0.1:9191

Starting third server and join the cluster
./yadosctl server start --name server3 --listen-address 127.0.0.1 --wal-dir /tmp --port 9193 --log-level info --peer server1:127.0.0.1:9191 --peer server2:127.0.0.1:9192

Usage:
  yadosctl server start [flags]

Flags:
  -h, --help                    help for start
      --http-port int           Port to use for http server (default -1)
      --listen-address string   Listen Address on which server will listen
                                Can usually be left blank. Otherwise, use IP address or host Name 
                                that other server nodes use to connect to the new server (default "127.0.0.1")
      --log-level string        Log level [info|debug|warn|trace|error] (default "info")
      --name string             Name of the server, Name must be unique in a cluster
      --peer strings            peer to join <name:ip-address:port>, use multiple of this flag if want to join with multiple peers
      --port int32              Port to use for communication (default 9191)
      --wal-dir string          Location for replicated write ahead log (default "info")
```
### Listing the members in the cluster
```shell
./yadosctl  server list --server 127.0.0.1 --port 9191
```
### Creating a store
```shell
./yadosctl store create --name store1
./yadosctl store create --name store1 --type kv
./yadosctl store create --name store1 --type sqlite
```
### Listing stores
```shell
./yadosctl store list
```
### Insert Object in Store
```shell
./yadosctl store put --key key3 --value value3 --store-name store1
```
### Get Object from a store
```shell
./yadosctl store get --store-name store1 --key key3
```

## Developer Notes
### How to build
```shell
1. Checkout the repository
2. Building using following command
   2.1 make
3. Run the tests using following command
   3.1 make test
4. Running single Test 
   4.1 TEST_NAME=<TestName> make test-single
```

### Fixing some common lint errors
```shell
Error      : File is not `gofmt`-ed with `-s
Resolution : Run the gofmt tool on the file which is reported this error

Error      : File is not `goimports`-ed (goimports)
Resolution : Run following command for the file

$GOPATH/bin/goimports  -w <file-name> 
```