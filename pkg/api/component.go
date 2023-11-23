package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateComponent(data *model.Component) (*model.Component, error) {
	fmt.Println("Create component api handler")
	return data, nil
}

func (t *TitondAPI) GetComponentByType(networkID uint, componentType string) (*model.Component, error) {
	fmt.Println("Get Component By Type:", componentType, " | NetworkID:", networkID)
	return &model.Component{Name: "GetComponentByType"}, nil
}

func (t *TitondAPI) GetComponentById(componentID uint) (*model.Component, error) {
	fmt.Println("Create component by id:", componentID)
	return &model.Component{Name: "GetComponentById"}, nil
}
