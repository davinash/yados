package proto

//go:generate protoc -I . -I /home/adongre/tools/protobuf-3.17.3/src --go_out=. --go-grpc_out=. common.proto
//go:generate protoc -I . -I /home/adongre/tools/protobuf-3.17.3/src --go_out=. --go-grpc_out=. service.proto
//go:generate protoc -I . -I /home/adongre/tools/protobuf-3.17.3/src --go_out=. --go-grpc_out=. controller.proto
