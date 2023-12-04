package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
	"gorm.io/gorm"
)

type MockTitondAPI struct {
	network   *model.Network
	networks  []model.Network
	component *model.Component
	err       error
	page      int
	networkID uint
}

func (mock *MockTitondAPI) CreateNetwork(data *model.Network) (*model.Network, error) {
	return mock.network, mock.err
}

func (mock *MockTitondAPI) GetNetworksByPage(page int) ([]model.Network, error) {
	mock.page = page
	return mock.networks, mock.err
}

func (mock *MockTitondAPI) GetNetworkByID(networkID uint) (*model.Network, error) {
	mock.networkID = networkID
	return mock.network, mock.err
}

func (mock *MockTitondAPI) DeleteNetwork(networkID uint) error {
	mock.networkID = networkID
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
		t.Fatal(" ---------- ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.R.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetNetworkByPage(t *testing.T) {
	testcases := []struct {
		link             string
		mockNetworks     []model.Network
		mockError        error
		expectedHttpCode int
		expectedPage     int
	}{
		{
			link:             "/api/networks/?page=1",
			mockNetworks:     []model.Network{},
			mockError:        types.ErrResourceNotFound,
			expectedHttpCode: http.StatusNotFound,
			expectedPage:     1,
		},
		{
			link:             "/api/networks/",
			mockNetworks:     []model.Network{},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
			expectedPage:     1,
		},
		{
			link:             "/api/networks/?page=abc",
			mockNetworks:     []model.Network{},
			mockError:        nil,
			expectedHttpCode: http.StatusBadRequest,
			expectedPage:     0,
		},
		{
			link:             "/api/networks/?page=1",
			mockNetworks:     []model.Network{},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
			expectedPage:     1,
		},
		{
			link:             "/api/networks/?page=0",
			mockNetworks:     []model.Network{},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
			expectedPage:     1,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		titondAPI.err = testcase.mockError
		titondAPI.networks = testcase.mockNetworks
		req, err := http.NewRequest(http.MethodGet, testcase.link, bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)
		assert.Equal(t, testcase.expectedHttpCode, w.Code)
		if testcase.expectedHttpCode != http.StatusBadRequest {
			assert.Equal(t, testcase.expectedPage, titondAPI.page)
		}
	}

}

func TestGetNetworkByID(t *testing.T) {
	testcases := []struct {
		link              string
		mockNetwork       *model.Network
		mockError         error
		expectedHttpCode  int
		expectedNetworkID uint
	}{
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         types.ErrResourceNotFound,
			expectedHttpCode:  http.StatusNotFound,
			expectedNetworkID: 12,
		},
		{
			link:              "/api/networks/1a",
			mockNetwork:       &model.Network{},
			mockError:         nil,
			expectedHttpCode:  http.StatusBadRequest,
			expectedNetworkID: 1,
		},
		{
			link:              "/api/networks/1",
			mockNetwork:       &model.Network{},
			mockError:         gorm.ErrRecordNotFound,
			expectedHttpCode:  http.StatusNotFound,
			expectedNetworkID: 1,
		},
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         nil,
			expectedHttpCode:  http.StatusOK,
			expectedNetworkID: 12,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		titondAPI.err = testcase.mockError
		titondAPI.network = testcase.mockNetwork
		req, err := http.NewRequest(http.MethodGet, testcase.link, bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)
		assert.Equal(t, testcase.expectedHttpCode, w.Code)
		if testcase.expectedHttpCode != http.StatusBadRequest {
			assert.Equal(t, testcase.expectedNetworkID, titondAPI.networkID)
		}
	}
}

func TestDeleteNetworkByID(t *testing.T) {
	testcases := []struct {
		link              string
		mockNetwork       *model.Network
		mockError         error
		expectedHttpCode  int
		expectedNetworkID uint
	}{
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         types.ErrResourceNotFound,
			expectedHttpCode:  http.StatusNotFound,
			expectedNetworkID: 12,
		},
		{
			link:              "/api/networks/12",
			mockNetwork:       nil,
			mockError:         types.ErrResourceNotFound,
			expectedHttpCode:  http.StatusNotFound,
			expectedNetworkID: 12,
		},
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         types.ErrInternalServer,
			expectedHttpCode:  http.StatusInternalServerError,
			expectedNetworkID: 12,
		},
		{
			link:              "/api/networks/1a",
			mockNetwork:       &model.Network{},
			mockError:         nil,
			expectedHttpCode:  http.StatusBadRequest,
			expectedNetworkID: 1,
		},
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         nil,
			expectedHttpCode:  http.StatusOK,
			expectedNetworkID: 12,
		},
		{
			link:              "/api/networks/12",
			mockNetwork:       &model.Network{},
			mockError:         errors.New("unknown error will return code 500"),
			expectedHttpCode:  http.StatusInternalServerError,
			expectedNetworkID: 12,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		titondAPI.err = testcase.mockError
		titondAPI.network = testcase.mockNetwork
		req, err := http.NewRequest(http.MethodDelete, testcase.link, bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)
		assert.Equal(t, testcase.expectedHttpCode, w.Code)
		if testcase.expectedHttpCode != http.StatusBadRequest {
			assert.Equal(t, testcase.expectedNetworkID, titondAPI.networkID)
		}
	}
}

func TestCreateComponent(t *testing.T) {
	testcases := []struct {
		link             string
		body             *model.Component
		mockError        error
		expectedHttpCode int
	}{
		{
			link:             "/api/components/",
			body:             &model.Component{},
			mockError:        types.ErrBadRequest,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link: "/api/components/",
			body: &model.Component{
				Name: "Titan-test",
				Type: "l2geth",
			},
			mockError:        types.ErrBadRequest,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link: "/api/components/",
			body: &model.Component{
				Name: "Titan-test",
				Type: "l2-geth",
			},
			mockError:        types.ErrInvalidComponentType,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link: "/api/components/",
			body: &model.Component{
				Name:      "Titan-test",
				Type:      "l2geth",
				NetworkID: 1,
			},
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link: "/api/components/",
			body: &model.Component{
				Name:      "Titan-test",
				Type:      "l2geth",
				NetworkID: 1,
			},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		fmt.Println(" Case: ", testcase)
		titondAPI.err = testcase.mockError
		bodyData, err := json.Marshal(testcase.body)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, testcase.link, bytes.NewBuffer(bodyData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)

		assert.Equal(t, testcase.expectedHttpCode, w.Code)
	}
}

func TestGetComponentByType(t *testing.T) {
	testcases := []struct {
		link             string
		mockComponent    *model.Component
		mockError        error
		expectedHttpCode int
	}{
		{
			link:             "/api/components/?type=l2geth&&network_id=108",
			mockComponent:    &model.Component{},
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link:             "/api/components/",
			mockComponent:    nil,
			mockError:        nil,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link:             "/api/components/?type=l2geth&&networkid=108",
			mockComponent:    &model.Component{},
			mockError:        types.ErrBadRequest,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link: "/api/components/?type=l2geth&&network_id=108",
			mockComponent: &model.Component{
				Name:      "Titan-test",
				Type:      "l2geth",
				NetworkID: 1,
			},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		fmt.Println(" Case: ", testcase)
		titondAPI.err = testcase.mockError
		titondAPI.component = testcase.mockComponent
		req, err := http.NewRequest(http.MethodGet, testcase.link, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)

		assert.Equal(t, testcase.expectedHttpCode, w.Code)
	}
}

func TestGetComponentByID(t *testing.T) {
	testcases := []struct {
		link             string
		mockComponent    *model.Component
		mockError        error
		expectedHttpCode int
	}{
		{
			link:             "/api/components/1a2",
			mockComponent:    &model.Component{},
			mockError:        nil,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link:             "/api/components/12",
			mockComponent:    nil,
			mockError:        types.ErrResourceNotFound,
			expectedHttpCode: http.StatusNotFound,
		},
		{
			link:             "/api/components/12",
			mockComponent:    &model.Component{},
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link:             "/api/components/12",
			mockComponent:    nil,
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link:             "/api/components/12",
			mockComponent:    &model.Component{},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		fmt.Println(" Case: ", testcase)
		titondAPI.err = testcase.mockError
		titondAPI.component = testcase.mockComponent
		req, err := http.NewRequest(http.MethodGet, testcase.link, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)

		assert.Equal(t, testcase.expectedHttpCode, w.Code)
	}
}

func TestDeleteComponentByID(t *testing.T) {
	testcases := []struct {
		link             string
		mockComponent    *model.Component
		mockError        error
		expectedHttpCode int
	}{
		{
			link:             "/api/components/1a2",
			mockComponent:    &model.Component{},
			mockError:        nil,
			expectedHttpCode: http.StatusBadRequest,
		},
		{
			link:             "/api/components/12",
			mockComponent:    nil,
			mockError:        types.ErrResourceNotFound,
			expectedHttpCode: http.StatusNotFound,
		},
		{
			link:             "/api/components/12",
			mockComponent:    &model.Component{},
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link:             "/api/components/12",
			mockComponent:    nil,
			mockError:        types.ErrInternalServer,
			expectedHttpCode: http.StatusInternalServerError,
		},
		{
			link:             "/api/components/12",
			mockComponent:    &model.Component{},
			mockError:        nil,
			expectedHttpCode: http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)

	titondAPI := &MockTitondAPI{}

	server := NewHTTPServer(nil, titondAPI)

	for _, testcase := range testcases {
		fmt.Println(" Case: ", testcase)
		titondAPI.err = testcase.mockError
		titondAPI.component = testcase.mockComponent
		req, err := http.NewRequest(http.MethodDelete, testcase.link, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)

		assert.Equal(t, testcase.expectedHttpCode, w.Code)
	}
}
