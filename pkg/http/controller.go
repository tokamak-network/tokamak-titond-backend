package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (s *HTTPServer) CreateNetwork(c *gin.Context) {
	result, err := s.apis.CreateNetwork(&model.Network{})
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func (s *HTTPServer) DeleteNetwork(c *gin.Context) {
	networkID := c.Param("id")
	fmt.Println("Request delete a network: ", networkID)
	id, err := strconv.ParseInt(networkID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("network id is invalid"))
		return
	}
	result, err := s.apis.DeleteNetwork(uint(id))
	if err == nil {
		if result > 0 {
			c.JSON(http.StatusOK, "Deleted")
		} else {
			c.JSON(http.StatusNotFound, "Not found")
		}
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func (s *HTTPServer) CreateComponent(c *gin.Context) {
	result, err := s.apis.CreateComponent(&model.Component{})
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func (s *HTTPServer) UpdateComponent(c *gin.Context) {
	result, err := s.apis.UpdateComponent(&model.Component{})
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}
