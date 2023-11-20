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
	ReadNetwork()
	ReadAllNetwork()
	UpdateNetwork()
	DeleteNetwork()
	CreateComponent()
	ReadComponent()
	ReadAllComponent()
	UpdateComponent()
	DeleteComponent()
}
