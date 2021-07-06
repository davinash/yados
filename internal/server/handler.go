package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	operationHandler map[OperationID]func(interface{}, *Server) (*Response, error)
)

func handleMessage(op Request, server *Server) (*Response, error) {
	if fn, ok := operationHandler[op.ID]; ok {
		response, err := fn(op.Arguments, server)
		if err != nil {
			response.Error = err.Error()
			return response, err
		}
		return response, nil
	}
	return &Response{
		ID:    "",
		Resp:  nil,
		Error: fmt.Errorf("unexpected Operation %v", op.ID).Error(),
	}, fmt.Errorf("unexpected Operation %v", op.ID)

}

func initialize() error {
	operationHandler = make(map[OperationID]func(interface{}, *Server) (*Response, error))

	operationHandler[PutObject] = PutObjectFn
	operationHandler[GetObject] = GetObjectFn
	operationHandler[DeleteObject] = DeleteObjectFn
	operationHandler[CreateStoreInCluster] = CreateStoreInClusterFn
	operationHandler[DeleteStoreFromCluster] = DeleteStoreFromClusterFn
	operationHandler[AddNewMember] = AddNewMemberFn
	operationHandler[StopServer] = StopServerFn
	operationHandler[AddNewMemberEx] = AddNewMemberExFn
	operationHandler[ListMembers] = ListMemberFn
	return nil
}

//SetupRouter Router setup function
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
