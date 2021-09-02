package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"

	pb "github.com/davinash/yados/internal/proto/gen"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type storeSqlite struct {
	db        *sql.DB
	logger    *logrus.Logger
	storeType pb.StoreType
	name      string
}

//NewSqliteStore creates a new store for sqlite database
func NewSqliteStore(args *Args) (SQLStore, error) {
	s := &storeSqlite{
		logger:    args.Logger,
		storeType: args.StoreType,
		name:      args.Name,
	}
	dbPath := filepath.Join(args.WALDir, fmt.Sprintf("%s.db", args.Name))
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	s.db = db
	return s, nil
}

func (ss *storeSqlite) Name() string {
	return ss.name
}

func (ss *storeSqlite) Close() error {
	return ss.db.Close()
}

func (ss *storeSqlite) Type() pb.StoreType {
	return ss.storeType
}

func (ss *storeSqlite) Execute(request *pb.ExecuteQueryRequest) (*pb.ExecuteQueryReply, error) {
	reply := &pb.ExecuteQueryReply{}

	result, err := ss.db.Exec(request.SqlQuery)
	if err != nil {
		reply.Error = err.Error()
		return reply, err
	}
	reply.RowsAffected, err = result.RowsAffected()
	if err != nil {
		reply.Error = err.Error()
		return reply, err
	}

	return reply, nil
}

func (ss *storeSqlite) Query(request *pb.QueryRequest) (*pb.QueryReply, error) {
	reply := &pb.QueryReply{}
	rows, err := ss.db.Query(request.SqlQuery)
	if err != nil {
		return reply, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ss.logger.Warnf("failed to close rows, Error = %v", err)
		}
	}(rows)

	columns, err := rows.Columns()
	if err != nil {
		return reply, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	reply.Rows = make([]*pb.TableRow, 0)

	for rows.Next() {
		tRow := pb.TableRow{
			Row: make([]*pb.ColumnValue, 0),
		}
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return reply, err
		}
		for i, col := range columns {
			val := values[i]
			bytes, err := json.Marshal(val)
			if err != nil {
				return reply, err
			}
			tRow.Row = append(tRow.Row, &pb.ColumnValue{
				Name:  col,
				Value: bytes,
			})
		}
		reply.Rows = append(reply.Rows, &tRow)
	}
	return reply, nil
}

func (ss *storeSqlite) Create(request *pb.StoreCreateRequest) {
	panic("implement me")
}

func (ss *storeSqlite) Delete(request *pb.StoreDeleteRequest) {
	panic("implement me")
}
