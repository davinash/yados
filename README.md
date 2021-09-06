# yados
<img src="https://github.com/davinash/yados/blob/main/logo.jpg" align="center"/>

Yet Another Distributed Object Store

[![Go Actions Status](https://github.com/davinash/yados/workflows/Go/badge.svg)](https://github.com/davinash/yados/actions)

## Command Line Interface
* [Server Commands](server.md)
* [Store Commands](store.md)

### Quick Start
Using Docker Compose is basically a two-step process:

```shell
wget https://raw.githubusercontent.com/davinash/yados/main/docker-compose.yaml
docker-compose up -d
docker ps -a 
Above command should show three container instances with name yados-1, yados-2 and yados-3 running
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