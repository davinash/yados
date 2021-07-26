package server

import "sync"

//Log represents the interface for the persistent log
type Log interface {
	CurrentIndex() uint64
}

type log struct {
	mutex sync.RWMutex
}

//NewLog creates a new instance of the internal log
func NewLog() Log {
	return &log{}
}

func (l *log) CurrentIndex() uint64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return 0
}
