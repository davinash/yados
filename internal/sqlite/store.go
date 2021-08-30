package sqlite

import (
	"database/sql"
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

func (ss *storeSqlite) ExecuteDDLQuery(request *pb.DDLQueryRequest) (*pb.DDLQueryReply, error) {

	rows, err := ss.db.Query(request.SqlQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ss.logger.Warnf("failed to close the rows handle, Error = %v", err)
		}
	}(rows)

	ss.logger.Debug(rows)

	return &pb.DDLQueryReply{}, nil
}
