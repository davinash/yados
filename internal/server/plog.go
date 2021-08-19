package server

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/protobuf/types/known/anypb"

	"google.golang.org/protobuf/proto"

	pb "github.com/davinash/yados/internal/proto/gen"

	"github.com/sirupsen/logrus"
)

//PersistentLogFile name of the log file
var PersistentLogFile = "pLog.bin"

//PLog represents the pLog
type PLog interface {
	Open() error
	Close() error
	Append(entry *pb.LogEntry) error
	Size() int
	PLogFileName() string
	Iterator() (PLogIterator, error)
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
	mutex         sync.RWMutex
	logDir        string
	logger        *logrus.Entry
	storeFileName string
	storeFh       *os.File
	server        Server
	size          int
}

func (m *plog) PLogFileName() string {
	return m.storeFileName
}

//Open Open the pLog
func (m *plog) Open() error {
	m.storeFileName = filepath.Join(m.logDir, PersistentLogFile)
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

	buffer, err := proto.Marshal(entry)
	if err != nil {
		return err
	}
	var buf [binary.MaxVarintLen32]byte
	encodedLength := binary.PutUvarint(buf[:], uint64(len(buffer)))
	_, err = m.storeFh.Write(buf[:encodedLength])
	if err != nil {
		return err
	}
	_, err = m.storeFh.Write(buffer)
	if err != nil {
		return err
	}
	m.size++

	if m.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		var commandStr []byte

		switch entry.CmdType {
		case pb.CommandType_CreateStore:
			var command pb.StoreCreateRequest
			err = anypb.UnmarshalTo(entry.Command, &command, proto.UnmarshalOptions{})
			if err != nil {
			}
			commandStr, err = json.Marshal(&command)
			if err != nil {
			}

		case pb.CommandType_Put:
			var command pb.PutRequest
			err = anypb.UnmarshalTo(entry.Command, &command, proto.UnmarshalOptions{})
			if err != nil {
			}
			commandStr, err = json.Marshal(&command)
			if err != nil {
			}
		}
		m.logger.Debugf("[%s] Entry Appended (Term = %v, Index = %v ) Value = %s",
			entry.Id, entry.Term, entry.Index, string(commandStr))
	}

	if m.server.EventHandler() != nil {
		if m.size == m.server.EventHandler().PersistEntryEventThreshold() {
			m.server.EventHandler().SendEvent(m.size)
		}
	}
	return nil
}

func (m *plog) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.size
}

//PLogIterator iterator interface for persistent logs
type PLogIterator interface {
	Next() (*pb.LogEntry, error)
	Close() error
}

type pLogIterator struct {
	pLog    PLog
	storeRH *os.File
}

//Iterator creates a new instance of the iterator
func (m *plog) Iterator() (PLogIterator, error) {
	iter := pLogIterator{
		pLog: m,
	}
	storeRH, err := os.Open(m.PLogFileName())
	if err != nil {
		return nil, err
	}
	iter.storeRH = storeRH
	return &iter, nil
}

//ErrInvalidVarint error when the data is corrupt
var ErrInvalidVarint = errors.New("invalid varint32 encountered")

func (p *pLogIterator) Next() (*pb.LogEntry, error) {

	var headerBuf [binary.MaxVarintLen32]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 {
		if bytesRead >= len(headerBuf) {
			return nil, ErrInvalidVarint
		}

		newBytesRead, err := p.storeRH.Read(headerBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if err != nil {
				if err == io.EOF {
					return nil, nil
				}
				return nil, err
			}
			continue
		}
		bytesRead += newBytesRead
		messageLength, varIntBytes = decodeVarint(headerBuf[:bytesRead])
	}
	messageBuf := make([]byte, messageLength)
	_, err := io.ReadFull(p.storeRH, messageBuf)
	if err != nil {
		return nil, err
	}
	entry := pb.LogEntry{}
	err = proto.Unmarshal(messageBuf, &entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func decodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}
	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

func (p *pLogIterator) Close() error {
	err := p.storeRH.Close()
	if err != nil {
		return err
	}
	return nil
}
