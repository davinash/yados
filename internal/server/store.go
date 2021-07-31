package server

import "github.com/sirupsen/logrus"

//Store represents the store
type Store interface {
	Open() error
	Close() error
	Put() error
	Get() error
}

//NewStorage Creates new storage
func NewStorage(logDir string, logger *logrus.Entry) (Store, error) {
	ms := &MemoryStore{
		logDir: logDir,
		logger: logger,
	}
	return ms, nil
}
