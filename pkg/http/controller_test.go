package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

type MockTitondAPI struct {
	network   *model.Network
	networks  []model.Network
	component *model.Component
	err       error
}

func (mock *MockTitondAPI) CreateNetwork(data *model.Network) (*model.Network, error) {
	fmt.Println("Create Network in Mock")
	return mock.network, mock.err
}

func (mock *MockTitondAPI) GetNetworksByPage(page int) ([]model.Network, error) {
	return mock.networks, mock.err
}

func (mock *MockTitondAPI) GetNetworkByID(networkID uint) (*model.Network, error) {
	return mock.network, mock.err
}

func (mock *MockTitondAPI) DeleteNetwork(id uint) error {
	return mock.err
}

func (mock *MockTitondAPI) CreateComponent(component *model.Component) (*model.Component, error) {
	return mock.component, mock.err
}

func (mock *MockTitondAPI) GetComponentByType(networkID uint, componentType string) (*model.Component, error) {
	return mock.component, mock.err
}

func (mock *MockTitondAPI) GetComponentById(componentID uint) (*model.Component, error) {
	return mock.component, mock.err
}

func (mock *MockTitondAPI) DeleteComponentById(componentID uint) error {
	return mock.err
}

func TestCreateNetworkSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}
	server := NewHTTPServer(nil, titondAPI)

	req, err := http.NewRequest(http.MethodPost, "/api/networks/", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.R.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	titondAPI.network = nil
	titondAPI.err = types.ErrInternalServer

}

func TestCreateNetworkFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}
	titondAPI.network = nil
	titondAPI.err = types.ErrInternalServer

	server := NewHTTPServer(nil, titondAPI)

	req, err := http.NewRequest(http.MethodPost, "/api/networks/", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.R.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
