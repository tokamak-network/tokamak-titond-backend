package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func (t *TitondAPI) CreateComponent(component *model.Component) (*model.Component, error) {
	var result *model.Component
	var err error

	if err := t.checkNetworkStatus(component.NetworkID); err != nil {
		return nil, err
	}

	namespace := generateNamespace(component.NetworkID)

	if t.checkNamespace(namespace) != nil {
		return nil, err
	}

	switch component.Type {
	case "data-transport-layer":
		result, err = t.CreateDTL(component)

	case "l2geth":
		result, err = t.CreateL2Geth(component)

	case "batch-submitter":
		result, err = t.CreateBatchSubmitter(component)

	case "relayer":
		result, err = t.CreateRelayer(component)

	case "explorer":
		result, err = t.CreateExplorer(component)

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
	return t.db.ReadComponentByType(componentType, networkID)
}

func (t *TitondAPI) GetComponentById(componentID uint) (*model.Component, error) {
	fmt.Println("Get component by id:", componentID)
	return t.db.ReadComponent(componentID)
}

func (t *TitondAPI) DeleteComponentById(componentID uint) error {
	fmt.Println("Delete component by id:", componentID)
	result, err := t.db.DeleteComponent(componentID)
	if err == nil {
		if result == 0 {
			return types.ErrResourceNotFound
		}
	}
	return err
}

func (t *TitondAPI) checkNetworkStatus(networkID uint) error {
	network, err := t.db.ReadNetwork(networkID)
	if err != nil {
		return err
	}
	if !network.Status {
		return apptypes.ErrNetworkNotReady
	}
	return nil
}

func (t *TitondAPI) checkNamespace(namespace string) error {
	_, err := t.k8s.GetNamespace(namespace)
	return err
}
