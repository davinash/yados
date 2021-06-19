package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func handleMessage(op Operation, server *Server) (interface{}, error) {
	return nil, nil
}

func setupRouter(server *Server) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/message", func(context *gin.Context) {
		var op Operation
		if err := context.ShouldBindWith(&op, binding.JSON); err != nil {
			context.JSON(http.StatusBadRequest, err.Error())
			return
		}
		result, err := handleMessage(op, server)
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
	return router
}
