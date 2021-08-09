package server

//CommandType state of the raft instance
type CommandType int

const (
	//CreateStore Command Type for creating a new store
	CreateStore CommandType = iota
)

//Command command interface
type Command interface {
	Execute() error
	Type() CommandType
	Server() Server
}

type command struct {
	cmdType CommandType
	srv     Server
	args    interface{}
}

//NewCommand creates a new object of command
func NewCommand(srv Server, cmdType CommandType, args interface{}) Command {
	c := &command{
		srv:     srv,
		cmdType: cmdType,
		args:    args,
	}
	return c
}

func (c *command) Server() Server {
	return c.srv
}

func (c *command) Type() CommandType {
	return c.cmdType
}

func (c *command) Execute() error {
	switch c.Type() {
	case CreateStore:

	}
	return nil
}
