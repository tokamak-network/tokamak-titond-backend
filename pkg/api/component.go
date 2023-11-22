package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateComponent(data *model.Component) (*model.Component, error) {
	fmt.Println("Create component api handler")
	return data, nil
}

func (t *TitondAPI) UpdateComponent(data *model.Component) (*model.Component, error) {
	fmt.Println("Update component api handler")
	return data, nil
}
