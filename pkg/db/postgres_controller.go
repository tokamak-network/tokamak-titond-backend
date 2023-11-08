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

func (p *Postgres) ReadNetwork() {

}
func (p *Postgres) ReadAllNetwork() {

}
func (p *Postgres) UpdateNetwork() {

}
func (p *Postgres) DeleteNetwork() {

}
func (p *Postgres) CreateResource() {

}
func (p *Postgres) ReadResource() {

}
func (p *Postgres) ReadAllResource() {

}
func (p *Postgres) UpdateResource() {

}
func (p *Postgres) DeleteResource() {

}
