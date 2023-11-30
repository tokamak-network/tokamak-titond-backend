package db

import (
	"fmt"

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
	fmt.Println("Offset ", offset, " Limit: ", limit)
	var networks []model.Network
	result := p.gDB.Offset(offset).Limit(limit).Find(&networks)
	return networks, result.Error
}

func (p *Postgres) ReadAllNetwork() ([]model.Network, error) {
	var networks []model.Network
	result := p.gDB.Find(networks)
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
