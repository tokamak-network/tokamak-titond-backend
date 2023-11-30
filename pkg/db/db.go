package db

import "github.com/tokamak-network/tokamak-titond-backend/pkg/model"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Client interface {
	CreateNetwork(*model.Network) (*model.Network, error)
	ReadNetwork(uint) (*model.Network, error)
	ReadNetworkByRange(int, int) ([]model.Network, error)
	ReadAllNetwork() ([]model.Network, error)
	UpdateNetwork(network *model.Network) (*model.Network, error)
	DeleteNetwork(uint) (int64, error)
	CreateComponent()
	ReadComponent()
	ReadAllComponent()
	UpdateComponent()
	DeleteComponent()
}
