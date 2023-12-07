package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func (t *TitondAPI) CreateComponent(component *model.Component) (*model.Component, error) {
	var result *model.Component

	namespace := generateNamespace(component.NetworkID)
	_, err := t.k8s.GetNamespace(namespace)
	if err != nil {
		t.k8s.CreateNamespace(namespace)
	}

	switch component.Type {
	case "l2geth":
		result, err = t.CreateL2Geth(component)

	case "data-transport-layer":
		result, err = t.CreateDTL(component)

	case "batch-submitter":
		result, err = t.CreateBatchSubmitter(component)

	default:
		err = apptypes.ErrInvalidComponentType
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TitondAPI) GetComponentByType(networkID uint, componentType string) (*model.Component, error) {
	fmt.Println("Get Component By Type:", componentType, " | NetworkID:", networkID)
	return &model.Component{Name: "GetComponentByType"}, nil
}

func (t *TitondAPI) GetComponentById(componentID uint) (*model.Component, error) {
	fmt.Println("Get component by id:", componentID)
	return &model.Component{Name: "GetComponentById"}, nil
}

func (t *TitondAPI) DeleteComponentById(componentID uint) error {
	fmt.Println("Delete component by id:", componentID)
	return nil
}
