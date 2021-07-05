package server

import "syscall"

//Stop stops the members in the cluster
func Stop(args interface{}, server *Server) (*Response, error) {
	server.OSSignalCh <- syscall.SIGINT
	return nil, nil
}
