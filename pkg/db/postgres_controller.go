package db

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (p *Postgres) CreateNetwork(data *model.Network) (*model.Network, error) {
	result := p.gDB.Create(data)
	return data, result.Error
}

func (p *Postgres) ReadNetwork(networkID uint) (*model.Network, error) {
	var network model.Network
	result := p.gDB.First(&network, networkID)
	return &network, result.Error
}

func (p *Postgres) ReadAllNetwork() {

}

func (p *Postgres) UpdateNetwork(data *model.Network) (*model.Network, error) {
	result := p.gDB.Save(data)
	return data, result.Error
}

func (p *Postgres) DeleteNetwork(networkID uint) (int64, error) {
	result := p.gDB.Delete(&model.Network{}, networkID)

	return result.RowsAffected, result.Error
}

func (p *Postgres) CreateComponent(component *model.Component) (*model.Component, error) {
	result := p.gDB.Create(component)
	if result.Error != nil {
		return nil, result.Error
	}
	return component, nil
}
func (p *Postgres) ReadComponent() {

}
func (p *Postgres) ReadAllComponent() {

}
func (p *Postgres) UpdateComponent() {

}
func (p *Postgres) DeleteComponent() {

}
