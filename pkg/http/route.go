package http

import "github.com/gin-gonic/gin"

func (s *HTTPServer) NewTestRouter(group *gin.RouterGroup) {
	group.GET("/", s.TestAPI)
}

func (s *HTTPServer) NewNetworkRouter(group *gin.RouterGroup) {
	group.POST("/", s.CreateNetwork)
}

func (s *HTTPServer) NewResourceRouter(group *gin.RouterGroup) {

}
