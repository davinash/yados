package server

import "syscall"

func Stop(args interface{}, server *Server) (*Response, error) {
	server.OSSignalCh <- syscall.SIGINT
	return nil, nil
}
