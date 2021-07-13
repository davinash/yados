package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

//Store represents the store in the cluster
type Store struct {
}

//CreateStore Creates new store in the cluster
func (server *YadosServer) CreateStore(ctx context.Context, request *pb.StoreCreateRequest) (*pb.StoreCreateReply, error) {
	panic("implement me")
}
