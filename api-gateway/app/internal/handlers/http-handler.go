package handlers

import "github.com/gin-gonic/gin"

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("calculate")
		tasks := api.Group("/tasks")
		{
			tasks.POST("/")
			tasks.POST("/status/:id")
		}
	}

	return router
}
