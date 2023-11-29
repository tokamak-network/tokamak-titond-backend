package db

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (p *Postgres) CreateNetwork(data *model.Network) (*model.Network, error) {
	result := p.gDB.Create(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (p *Postgres) ReadNetwork(networkID uint) (*model.Network, error) {
	network := model.Network{}
	result := p.gDB.First(&network, networkID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &network, nil
}

func (p *Postgres) ReadAllNetwork() {

}

func (p *Postgres) UpdateNetwork(data *model.Network) (*model.Network, error) {
	result := p.gDB.Save(data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (p *Postgres) DeleteNetwork(networkID uint) (int64, error) {
	result := p.gDB.Delete(&model.Network{}, networkID)

	return result.RowsAffected, result.Error
}
func (p *Postgres) CreateComponent() {

}
func (p *Postgres) ReadComponent() {

}
func (p *Postgres) ReadAllComponent() {

}
func (p *Postgres) UpdateComponent() {

}
func (p *Postgres) DeleteComponent() {

}
