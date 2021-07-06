package server

import "syscall"

//StopServerFn stops the members in the cluster
func StopServerFn(args interface{}, server *Server) (*Response, error) {
	server.OSSignalCh <- syscall.SIGINT
	return nil, nil
}
