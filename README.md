[![Go Actions Status](https://github.com/davinash/yados/workflows/Go/badge.svg)](https://github.com/davinash/yados/actions)
# yados
Yet Another Distributed Object Store

## Starting a server
### Starting with default options
In this mode server will start listening on 127.0.0.1 and port 9191
```shell
./yados server start --name server1
```
Starting server on specific listen address and port
```shell
./yados server start --name server1 --listen-address 127.0.0.1 --port 9191
```
Starting a server and joining the cluster with other peer
```shell
/yados server start --listen-address 127.0.0.1 --port 9192 --name Server2 --peer 127.0.0.1:9191
```
Starting a server and joining the cluster with multiple peers
```shell
./yados server start --listen-address 127.0.0.1 --port 9193 --name Server3 --peer 127.0.0.1:9191 --peer 127.0.0.1:9192
```
### Listing the members in the cluster
```shell
./yadosctl  server list --server 127.0.0.1 --port 9191
```

### Stopping individual server in a cluster
```shell
./yadosctl  server stop --server 127.0.0.1 --port 9191
```

## Developer Notes
### How to build
```shell
1. Checkout the repository
2. Building using following command
   2.1 make
3. Run the tests using following command
   3.1 make test
```