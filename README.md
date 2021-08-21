# yados
<img src="https://github.com/davinash/yados/blob/main/logo.jpg" align="center"/>

Yet Another Distributed Object Store

[![Go Actions Status](https://github.com/davinash/yados/workflows/Go/badge.svg)](https://github.com/davinash/yados/actions)

## Starting a server
### Starting with default options
In this mode server will start listening on 127.0.0.1 and port 9191
```shell
./yadosctl server start --name server1
```
Starting server on specific listen address and port
```shell
./yadosctl server start --name server1 --listen-address 127.0.0.1 --port 9191
```
Starting a server and joining the cluster with other peer
```shell
./yadosctl server start --listen-address 127.0.0.1 --port 9192 --name server2 --peer server1:127.0.0.1:9191
```
Starting a server and joining the cluster with multiple peers
```shell
./yadosctl server start --listen-address 127.0.0.1 --port 9193 --name server3 --peer server1:127.0.0.1:9191 --peer server2:127.0.0.1:9192
```
### Listing the members in the cluster
```shell
./yadosctl  server list --server 127.0.0.1 --port 9191
```
### Creating a store
```shell
./yadosctl store create --name store1
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