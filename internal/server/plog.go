package server

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

//PLog represents the pLog
type PLog interface {
	Open() error
	Close() error
	Put() error
	Get() error
}

//NewPLog Creates new storage
func NewPLog(logDir string, logger *logrus.Entry) (PLog, error) {
	ms := &plog{
		logDir: logDir,
		logger: logger,
	}
	return ms, nil
}

//plog represents temporary memory store
type plog struct {
	logDir        string
	logger        *logrus.Entry
	storeFileName string
	storeFh       *os.File
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

//Put puts a new value in the pLog
func (m *plog) Put() error {
	panic("implement me")
}

//Get gets a value for the key
func (m *plog) Get() error {
	panic("implement me")
}
