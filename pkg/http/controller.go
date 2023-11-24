package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

// @Summary CreateNetwork
// @Description Create a new network
// @Produce json
// @Success 200 {object} model.Network
// @Router /api/networks [post]
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
		message := fmt.Sprintf("Deleted network id: %d", networkID)
		c.JSON(http.StatusOK, message)
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

func (s *HTTPServer) GetComponentByType(c *gin.Context) {
	var params model.Component
	if err := c.ShouldBindQuery(&params); err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	result, err := s.apis.GetComponentByType(params.NetworkID, params.Type)
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

func (s *HTTPServer) GetComponentById(c *gin.Context) {
	componentID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	result, err := s.apis.GetComponentById(uint(componentID))
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

func (s *HTTPServer) DeleteComponentById(c *gin.Context) {
	componentID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	err = s.apis.DeleteComponentById(uint(componentID))
	if err == nil {
		message := fmt.Sprintf("Deleted component %d", componentID)
		c.JSON(http.StatusOK, message)
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
