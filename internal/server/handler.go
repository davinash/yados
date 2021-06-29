package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
)

var (
	operationHandler map[OperationId]func(interface{}, *Server) (*Response, error)
)

func handleMessage(op Request, server *Server) (*Response, error) {
	if fn, ok := operationHandler[op.Id]; ok {
		return fn(op.Arguments, server)
	} else {

	}
	return &Response{}, nil
}

func initialize() error {
	operationHandler = make(map[OperationId]func(interface{}, *Server) (*Response, error))

	operationHandler[PutObject] = Put
	operationHandler[GetObject] = Get
	operationHandler[DeleteObject] = Delete
	operationHandler[CreateStoreInCluster] = CreateStore
	operationHandler[DeleteStoreFromCluster] = DeleteStore
	operationHandler[AddNewMember] = Join
	return nil
}

func SetupRouter(server *Server) *gin.Engine {
	initialize()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"response": "pong",
		})
	})

	router.POST("/message", func(context *gin.Context) {
		var op Request
		if err := context.ShouldBindWith(&op, binding.JSON); err != nil {
			log.Println(err.Error())
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}
		resp, err := handleMessage(op, server)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err)
			return
		}
		context.JSON(http.StatusOK, resp)
	})
	return router
}
