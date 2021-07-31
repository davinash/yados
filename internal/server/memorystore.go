package server

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

//MemoryStore represents temporary memory store
type MemoryStore struct {
	logDir        string
	logger        *logrus.Entry
	storeFileName string
	storeFh       *os.File
}

//Open Open the store
func (m *MemoryStore) Open() error {
	m.storeFileName = filepath.Join(m.logDir, "store.bin")
	file, err := os.OpenFile(m.storeFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	m.storeFh = file
	return nil
}

//Close close the store
func (m *MemoryStore) Close() error {
	err := m.storeFh.Close()
	if err != nil {
		return err
	}
	return nil
}

//Put puts a new value in the store
func (m *MemoryStore) Put() error {
	panic("implement me")
}

//Get gets a value for the key
func (m *MemoryStore) Get() error {
	panic("implement me")
}
