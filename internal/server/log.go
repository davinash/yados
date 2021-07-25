package server

//Log represents the interface for the persistent log
type Log interface {
}

type log struct {
}

//NewLog creates a new instance of the internal log
func NewLog() Log {
	return &log{}
}
