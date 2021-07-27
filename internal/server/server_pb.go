package server

import (
	"context"

	pb "github.com/davinash/yados/internal/proto/gen"
)

func (srv *server) RequestVotes(ctx context.Context, request *pb.VoteRequest) (*pb.VoteReply, error) {
	panic("implement me")
}

func (srv *server) AppendEntries(ctx context.Context, request *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	panic("implement me")
}
