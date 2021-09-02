package fsm

import (
	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/store"
)

//FSM interface for state machine to apply
type FSM interface {
	Apply(entry *pb.LogEntry) error
	Store() store.Store
}

type fsm struct {
	store store.Store
}

//NewFSM creates new instance of fsm
func NewFSM(store store.Store) FSM {
	f := fsm{}
	return &f
}

func (f *fsm) Store() store.Store {
	return f.store
}

func (f *fsm) Apply(entry *pb.LogEntry) error {
	//switch entry.CmdType {
	//
	//case pb.CommandType_CreateStore:
	//	var req pb.StoreCreateRequest
	//	err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
	//	if err != nil {
	//		return err
	//	}
	//	err = f.Store().StoreCreate(&req)
	//	if err != nil {
	//		return err
	//	}
	//
	//case pb.CommandType_Put:
	//	var req pb.PutRequest
	//	err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
	//	if err != nil {
	//		return err
	//	}
	//	return (f.Store().Stores()[req.StoreName].(kv.KVStore)).Put(&req)
	//
	//case pb.CommandType_DeleteStore:
	//	var req pb.StoreDeleteRequest
	//	err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
	//	if err != nil {
	//		return err
	//	}
	//	err = srv.StoreDelete(&req)
	//	if err != nil {
	//		return err
	//	}
	//
	//case pb.CommandType_SqlDDL:
	//	var req pb.ExecuteQueryRequest
	//	err := anypb.UnmarshalTo(entry.Command, &req, proto.UnmarshalOptions{})
	//	if err != nil {
	//		return err
	//	}
	//	q, err := (srv.Stores()[req.StoreName].(kv.SQLStore)).Execute(&req)
	//	if err != nil {
	//		return err
	//	}
	//	srv.logger.Debug(q)
	//}
	return nil
}
