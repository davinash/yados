package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
)

type Response struct {
	response interface{}
	err      error
}

var (
	operationHandler map[string]func(interface{}) Response
)

func handleMessage(op Operation, server *Server) Response {
	if fn, ok := operationHandler[op.Name]; ok {
		return fn(op.Arguments)
	} else {

	}
	return Response{}
}

func initialize() error {
	operationHandler = make(map[string]func(interface{}) Response)

	operationHandler["Put"] = Put
	operationHandler["Get"] = Get
	operationHandler["Delete"] = Delete
	operationHandler["CreateStore"] = CreateStore
	operationHandler["DeleteStore"] = DeleteStore
	return nil
}

func setupRouter(server *Server) *gin.Engine {
	initialize()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"response": "pong",
		})
	})

	router.POST("/message", func(context *gin.Context) {
		var op Operation
		if err := context.ShouldBindWith(&op, binding.JSON); err != nil {
			log.Println(err.Error())
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}
		resp := handleMessage(op, server)
		if resp.err != nil {
			context.JSON(http.StatusInternalServerError, resp.err.Error())
			return
		}
		resultMarshalled, err := json.Marshal(resp.response)
		if err != nil {
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, resultMarshalled)
	})
	return router
}
