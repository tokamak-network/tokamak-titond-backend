package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (s *HTTPServer) CreateNetwork(c *gin.Context) {
	result, err := s.apis.CreateNetwork(&model.Network{
		ContractAddressURL: "",
		StateDumpURL:       "",
		Status:             false,
	})
	// result, err := s.apis.CreateNetwork()
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}
