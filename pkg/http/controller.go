package http

import (
	"net/http"

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
