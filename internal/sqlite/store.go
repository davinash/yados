package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"

	pb "github.com/davinash/yados/internal/proto/gen"
	"github.com/davinash/yados/internal/store"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

type storeSqlite struct {
	db        *sql.DB
	logger    *logrus.Entry
	storeType pb.StoreType
	name      string
}

//NewSqliteStore creates a new store for sqlite database
func NewSqliteStore(args *store.Args) (store.SQLStore, error) {
	s := &storeSqlite{
		logger:    args.Logger,
		storeType: args.StoreType,
		name:      args.Name,
	}
	dbPath := filepath.Join(args.PLogDir, fmt.Sprintf("%s.db", args.Name))
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

func (ss *storeSqlite) Delete() error {
	panic("implement me")
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

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return reply, err
		}
		row := pb.SingleRow{
			Row: map[string][]byte{},
		}

		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			var v interface{}
			if ok {
				v = string(b)
			} else {
				v = val
			}
			bytes, err := json.Marshal(v)
			if err != nil {
				return reply, err
			}
			row.Row[col] = bytes
		}
		reply.AllRows = append(reply.AllRows, &row)
	}
	return reply, nil
}
