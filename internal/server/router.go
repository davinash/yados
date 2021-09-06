package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	//ignore
	swaggerFiles "github.com/swaggo/files"
	// ignore
	_ "github.com/davinash/yados/docs"
)

//HTTPHandler http handler
type HTTPHandler struct {
	srv Server
}

//NewHTTPHandler create object of the router
func NewHTTPHandler(srv Server) *HTTPHandler {
	return &HTTPHandler{
		srv: srv,
	}
}

//Router return the Gin router object
func (h *HTTPHandler) Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	url := ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", h.srv.Address(), h.srv.HTTPPort()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/api/v1/status", h.getStatus)
	router.POST("/api/v1/store", h.createStore)
	router.GET("/api/v1/stores", h.getStores)
	router.POST("/api/v1/kv/put", h.put)
	router.GET("/api/v1/kv/get/:storeName/:key", h.get)
	router.POST("/api/v1/sqlite/execute", h.execute)
	router.POST("/api/v1/sqlite/query", h.query)

	return router
}

//
// @Summary Get the status of the cluster
// @Description Get the status of the cluster
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /status [get]
func (h *HTTPHandler) getStatus(c *gin.Context) {
	status, err := ExecuteCmdStatus(h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

//
// @Summary Create new store in a cluster
// @Description Create new store in a cluster
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Param   CreateCommandArgs     body CreateCommandArgs     true        "CreateCommandArgs"
// @Router /store [post]
func (h *HTTPHandler) createStore(c *gin.Context) {
	request := CreateCommandArgs{}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		h.srv.Logger().Error(err)
		return
	}

	err := ExecuteCmdCreateStore(&request, h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//
// @Summary Get all the store in a cluster
// @Description Get all the store in a cluster
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Router /stores [get]
func (h *HTTPHandler) getStores(c *gin.Context) {
	status, err := ExecuteCmdListStore(h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

//
// @Summary Perform put operation in KV type of store
// @Description Perform put operation in KV type of store
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Param   PutArgs     body PutArgs     true        "PutArgs"
// @Router /kv/put [post]
func (h *HTTPHandler) put(c *gin.Context) {
	request := PutArgs{}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := ExecuteCmdPut(&request, h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//
// @Summary Gets value of a key
// @Description Gets value of a key
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Param   storeName     path    string     true        "storeName"
// @Param   key     path    string     true        "key"
// @Router /kv/get [get]
func (h *HTTPHandler) get(c *gin.Context) {
	request := GetArgs{}
	storeName := c.Params.ByName("storeName")
	key := c.Params.ByName("key")

	request.StoreName = storeName
	request.Key = key

	reply, err := ExecuteCmdGet(&request, h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, reply)
}

//
// @Summary Perform sql execute operation in sqlite type of store
// @Description Perform sql execute operation in sqlite type of store
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Param   QueryArgs     body QueryArgs     true        "QueryArgs"
// @Router /sqlite/execute [post]
func (h *HTTPHandler) execute(c *gin.Context) {
	request := QueryArgs{}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	reply, err := ExecuteCmdQuery(&request, h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, reply)
}

//
// @Summary Perform sql query operation in sqlite type of store
// @Description Perform sql query operation in sqlite type of store
// @Success 200 {string} string	"ok"
// @Accept  json
// @Produce  json
// @Param   QueryArgs     body QueryArgs     true        "QueryArgs"
// @Router /sqlite/query [post]
func (h *HTTPHandler) query(c *gin.Context) {
	request := QueryArgs{}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	reply, err := ExecuteCmdSQLQuery(&request, h.srv.Self().Address, h.srv.Self().Port)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, reply)
}
