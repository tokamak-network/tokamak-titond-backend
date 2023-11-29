package http

import "github.com/gin-gonic/gin"

func (s *HTTPServer) NewNetworkRouter(group *gin.RouterGroup) {
	group.POST("/", s.CreateNetwork)
	group.GET("/", s.GetNetworksByPage)
	group.GET("/:id", s.GetNetworkById)
	group.DELETE("/:id", s.DeleteNetwork)
}

func (s *HTTPServer) NewComponentRouter(group *gin.RouterGroup) {
	group.POST("/", s.CreateComponent)
	group.GET("/", s.GetComponentByType)
	group.GET("/:id", s.GetComponentById)
	group.DELETE("/:id", s.DeleteComponentById)
}
