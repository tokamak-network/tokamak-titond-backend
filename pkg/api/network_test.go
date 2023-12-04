package api

import (
	"testing"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

type MockDBClient struct {
	network      *model.Network
	networks     []model.Network
	numOfDeleted int64
	component    *model.Component
	err          error
	networkID    uint
	offset       int
	limit        int
}

func (client *MockDBClient) CreateNetwork(network *model.Network) (*model.Network, error) {
	return client.network, client.err
}

func (client *MockDBClient) ReadNetwork(networkID uint) (*model.Network, error) {
	client.networkID = networkID
	return client.network, client.err
}

func (client *MockDBClient) ReadNetworkByRange(offset int, limit int) ([]model.Network, error) {
	client.offset = offset
	client.limit = limit
	return client.networks[offset : offset+limit], client.err
}

func (client *MockDBClient) UpdateNetwork(network *model.Network) (*model.Network, error) {
	return client.network, client.err
}

func (client *MockDBClient) DeleteNetwork(networkID uint) (int64, error) {
	return client.numOfDeleted, client.err
}

func (client *MockDBClient) CreateComponent(component *model.Component) (*model.Component, error) {
	return client.component, client.err
}

func (client *MockDBClient) ReadComponent() {

}

func (client *MockDBClient) ReadAllComponent() {

}

func (client *MockDBClient) UpdateComponent() {

}

func (client *MockDBClient) DeleteComponent() {

}

func TestCreateDeployer(t *testing.T) {

}
