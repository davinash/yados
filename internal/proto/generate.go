package proto

//go:generate protoc -I . --go_out=. --go-grpc_out=. messages.proto
//go:generate protoc -I . --go_out=. --go-grpc_out=. yados.proto
