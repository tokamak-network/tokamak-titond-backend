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
	ReadNetwork(networkID uint) (*model.Network, error)
	ReadAllNetwork()
	UpdateNetwork(network *model.Network) (*model.Network, error)
	DeleteNetwork(uint) (int64, error)
	CreateComponent(component *model.Component) (*model.Component, error)
	ReadComponent()
	ReadAllComponent()
	UpdateComponent()
	DeleteComponent()
}
