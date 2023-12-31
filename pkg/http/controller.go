package http

import (
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

// @Summary CreateNetwork
// @Description Create a new network
// @ID create-network
// @Produce json
// @Success 200 {object} model.Network
// @Failure 500
// @Router /api/networks [post]
func (s *HTTPServer) CreateNetwork(c *gin.Context) {
	result, err := s.apis.CreateNetwork(&model.Network{})
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

// @Summary GetNetworksByPage
// @Description Get networks by page
// @ID get-networks-by-page
// @Produce json
// @Param page path int true "The page number. Defaults to 1 if page not provided."
// @Success 200 {object} object
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/networks [get]
func (s *HTTPServer) GetNetworksByPage(c *gin.Context) {
	pageRequest := c.Query("page")
	var page int
	var err error
	if pageRequest == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageRequest)
		if err != nil {
			s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
			return
		}
		if page < 1 {
			page = 1
		}
	}
	result, err := s.apis.GetNetworksByPage(page)
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

// @Summary GetNetworkById
// @Description Get a network by id
// @ID get-network-by-id
// @Produce json
// @Param id path int true "Network ID"
// @Success 200 {object} object
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/networks/{id} [get]
func (s *HTTPServer) GetNetworkById(c *gin.Context) {
	networkID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		s.ResponseErrorMessage(c, apptypes.ErrBadRequest)
		return
	}
	result, err := s.apis.GetNetworkByID(uint(networkID))
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		s.ResponseErrorMessage(c, err)
	}
}

// @Summary DeleteNetwork
// @Description Delete a network by id
// @ID delete-network
// @Produce json
// @Param id path int true "Network ID"
// @Success 200 {object} object
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/networks/{id} [delete]
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

// @Summary CreateComponent
// @Description Create a new component
// @ID create-component
// @Accept json
// @Produce json
// @Param input body model.Component true "Component data to create"
// @Success 200 {object} model.Component
// @Failure 400
// @Failure 500
// @Router /api/components [post]
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

// @Summary GetComponentByType
// @Description Get Component By Type
// @ID get-component-by-type
// @Param type query string true "Component type (e.g., l2geth)"
// @Param network_id query integer true "Network ID"
// @Produce json
// @Success 200 {object} object
// @Failure 400
// @Failure 500
// @Router /api/components [get]
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

// @Summary GetComponentByID
// @Description Get Component By ID
// @ID get-component-by-id
// @Param id path int true "Component ID"
// @Success 200 {object} object
// @Failure 400
// @Failure 500
// @Router /api/components/{id} [get]
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

// @Summary DeleteComponentById
// @Description Delete Component By Id
// @ID delete-component-by-id
// @Param id path int true "Component ID"
// @Success 200 {object} object
// @Failure 400
// @Failure 500
// @Router /api/components/{id} [delete]
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
	fmt.Println("Error: ", err)
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
	case apptypes.ErrInvalidComponentType:
		{
			c.JSON(http.StatusBadRequest, err)
		}
	case gorm.ErrRecordNotFound:
		{
			c.JSON(http.StatusNotFound, err)
		}
	case apptypes.ErrComponentDependency:
		{
			c.JSON(http.StatusPreconditionFailed, err)
		}
	case apptypes.ErrNetworkNotReady:
		{
			c.JSON(http.StatusPreconditionFailed, err)
		}
	default:
		{
			c.JSON(http.StatusInternalServerError, err)
		}
	}
}
