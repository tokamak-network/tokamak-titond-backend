package http

import "github.com/gin-gonic/gin"

func (s *HTTPServer) NewNetworkRouter(group *gin.RouterGroup) {
	group.POST("/", s.CreateNetwork)
	group.DELETE("/:id", s.DeleteNetwork)
}

func (s *HTTPServer) NewComponentRouter(group *gin.RouterGroup) {

}
