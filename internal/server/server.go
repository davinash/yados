package server

import (
	"encoding/json"
	"fmt"
	"github.com/davinash/yados/internal/message"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net"
	"net/http"
	"strconv"
)

type Server interface {
	Start() error
	Stop() error
	Status() error
}

type HttpServer struct {
	name     string
	port     int
	listener net.Listener
}

func GetExistingServer() (*HttpServer, error) {
	server := HttpServer{}
	return &server, nil
}

func CreateNewServer(name string, port string) (*HttpServer, error) {
	server := HttpServer{
		name: name,
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	server.port = p
	return &server, nil
}

func (server *HttpServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	server.listener = listener
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/message", func(context *gin.Context) {
		var op message.Operation
		if err := context.ShouldBindWith(&op, binding.JSON); err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}
		result, err := message.MessageHandler(op)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		resultMarshalled, err := json.Marshal(result)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, resultMarshalled)
	})
	http.Serve(server.listener, router)
	return nil
}

func (server *HttpServer) Stop() error {
	return nil
}

func (server *HttpServer) Status() error {
	return nil
}
