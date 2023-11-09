package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (s *HTTPServer) CreateNetwork(c *gin.Context) {
	result := s.apis.CreateNetwork(&model.Network{
		ContractAddressURL: "aa.com",
		StateDumpURL:       "bb.com",
		Status:             false,
	})

	c.JSON(http.StatusOK, result)
}

func (s *HTTPServer) TestAPI(c *gin.Context) {
	fmt.Println(" Received new request ")
	time.Sleep(time.Second * 30)
	fmt.Println(" Waited 30 seconds ")
	c.JSON(http.StatusOK, "")
}
