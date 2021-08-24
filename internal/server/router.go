package server

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
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

	router.GET("/api/v1/status", h.getStatus)
	router.POST("/api/v1/store", h.createStore)
	router.GET("/api/v1/stores", h.getStores)
	router.POST("/api/v1/put", h.put)
	router.GET("/api/v1/get/:storeName/:key", h.get)

	return router
}

func (h *HTTPHandler) getStatus(c *gin.Context) {
	status, err := ExecuteCmdStatus(&StatusArgs{
		Address: h.srv.Self().Address,
		Port:    h.srv.Self().Port,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *HTTPHandler) createStore(c *gin.Context) {
	request := CreateCommandArgs{
		Address: h.srv.Self().Address,
		Port:    h.srv.Self().Port,
	}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		h.srv.Logger().Error(err)
		return
	}

	err := ExecuteCmdCreateStore(&request)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *HTTPHandler) getStores(c *gin.Context) {
	status, err := ExecuteCmdListStore(&ListArgs{
		Address: h.srv.Self().Address,
		Port:    h.srv.Self().Port,
	})
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *HTTPHandler) put(c *gin.Context) {
	request := PutArgs{
		Address: h.srv.Self().Address,
		Port:    h.srv.Self().Port,
	}
	if err := c.ShouldBindWith(&request, binding.JSON); err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := ExecuteCmdPut(&request)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *HTTPHandler) get(c *gin.Context) {
	request := GetArgs{
		Address: h.srv.Self().Address,
		Port:    h.srv.Self().Port,
	}
	storeName := c.Params.ByName("storeName")
	key := c.Params.ByName("key")

	request.StoreName = storeName
	request.Key = key

	reply, err := ExecuteCmdGet(&request)
	if err != nil {
		h.srv.Logger().Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, reply)
}
