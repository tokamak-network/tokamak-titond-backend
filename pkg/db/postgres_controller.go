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

func (p *Postgres) ReadNetworkByRange(offset int, limit int) ([]model.Network, error) {
	var networks []model.Network
	result := p.gDB.Offset(offset).Limit(limit).Find(&networks)
	return networks, result.Error
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
	return component, result.Error
}

func (p *Postgres) ReadComponent(componentID uint) (*model.Component, error) {
	var component model.Component
	result := p.gDB.Where(&model.Component{ID: componentID}).First(&component)
	return &component, result.Error
}

func (p *Postgres) ReadComponentByType(typ string, networkID uint) (*model.Component, error) {
	var component model.Component
	result := p.gDB.Where(&model.Component{Type: typ, NetworkID: networkID}).First(&component)
	return &component, result.Error
}

func (p *Postgres) ReadAllComponent() {

}

func (p *Postgres) UpdateComponent(component *model.Component) (*model.Component, error) {
	result := p.gDB.Save(component)
	return component, result.Error
}

func (p *Postgres) DeleteComponent(componentID uint) (int64, error) {
	result := p.gDB.Delete(&model.Component{}, componentID)

	return result.RowsAffected, result.Error

}
