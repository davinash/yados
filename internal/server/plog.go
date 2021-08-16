package server

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/davinash/yados/internal/proto/gen"
	"google.golang.org/protobuf/proto"

	"github.com/sirupsen/logrus"
)

//PLog represents the pLog
type PLog interface {
	Open() error
	Close() error
	Append(entry *pb.LogEntry) error
	Get() error
	Size() int
}

//NewPLog Creates new storage
func NewPLog(srv Server) (PLog, error) {
	ms := &plog{
		logDir: srv.LogDir(),
		logger: srv.Logger(),
		server: srv,
	}
	return ms, nil
}

//plog represents temporary memory store
type plog struct {
	mutex         sync.Mutex
	logDir        string
	logger        *logrus.Entry
	storeFileName string
	storeFh       *os.File
	server        Server
	size          int
}

//Open Open the pLog
func (m *plog) Open() error {
	m.storeFileName = filepath.Join(m.logDir, "pLog.bin")
	file, err := os.OpenFile(m.storeFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	m.storeFh = file
	return nil
}

//Close close the pLog
func (m *plog) Close() error {
	err := m.storeFh.Close()
	if err != nil {
		return err
	}
	return nil
}

//Append Append a new value in the pLog
func (m *plog) Append(entry *pb.LogEntry) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entryBytes, err := proto.Marshal(entry)
	if err != nil {
		return err
	}
	length := len(entryBytes)
	_, err = m.storeFh.WriteString(fmt.Sprintf("%d", length))
	if err != nil {
		return err
	}
	_, err = m.storeFh.Write(entryBytes)
	if err != nil {
		return err
	}
	m.size++

	if m.server.EventHandler() != nil {
		if m.size == m.server.EventHandler().PersistEntryEventThreshold() {
			m.server.EventHandler().SendEvent(m.size)
		}
	}

	return nil
}

func (m *plog) Size() int {
	return m.size
}

//Get gets a value for the key
func (m *plog) Get() error {
	panic("implement me")
}
