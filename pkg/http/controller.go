package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
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
	networkID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	err = s.apis.DeleteNetwork(uint(networkID))
	if err == nil {
		c.JSON(http.StatusOK, "Deleted")
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

func (s *HTTPServer) CreateComponent(c *gin.Context) {
	var component model.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	result, err := s.apis.CreateComponent(&component)
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

func (s *HTTPServer) UpdateComponent(c *gin.Context) {
	var component model.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	result, err := s.apis.UpdateComponent(&component)
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

func (s *HTTPServer) ResponseErrorMessage(c *gin.Context, err error) {
	switch err {
	case apptypes.ErrBadRequest:
		{
			c.JSON(http.StatusBadRequest, err)
		}
	case apptypes.ErrResourceNotFound:
		{
			c.JSON(http.StatusNotFound, err)
		}
	case apptypes.ErrInternalServer:
		{
			c.JSON(http.StatusInternalServerError, err)
		}
	default:
		{
			c.JSON(http.StatusInternalServerError, err)
		}
	}
}
