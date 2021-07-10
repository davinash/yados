package proto

//go:generate protoc -I . --go_out=. --go-grpc_out=. member.proto
//go:generate protoc -I . --go_out=. --go-grpc_out=. yados.proto
