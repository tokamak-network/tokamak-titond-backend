package api

import (
	"fmt"

	"github.com/emicklei/go-restful/v3/log"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func (t *TitondAPI) CreateComponent(component *model.Component) (*model.Component, error) {
	var result *model.Component
	var err error

	if err := t.checkNetworkStatus(component.NetworkID); err != nil {
		return nil, err
	}

	namespace := generateNamespace(component.NetworkID)

	if !t.checkNamespace(namespace) {
		t.k8s.CreateNamespace(namespace)
		t.createAccounts(namespace)
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

func (t *TitondAPI) checkNamespace(namespace string) bool {
	_, err := t.k8s.GetNamespace(namespace)
	if err != nil {
		return false
	}
	return true
}

func (t *TitondAPI) createAccounts(namespace string) {
	sequencerKey, address := generateKey()
	log.Printf("created sequencer account: %s\n", address)

	proposerKey, address := generateKey()
	log.Printf("created proposer account: %s\n", address)

	relayerKey, address := generateKey()
	log.Printf("created relayer account: %s\n", address)

	signerKey, address := generateKey()
	log.Printf("created block signer account: %s\n", address)

	stringData := map[string]string{
		"BATCH_SUBMITTER_SEQUENCER_PRIVATE_KEY": sequencerKey,
		"BATCH_SUBMITTER_PROPOSER_PRIVATE_KEY":  proposerKey,
		"MESSAGE_RELAYER__L1_WALLET":            relayerKey,
		"BLOCK_SIGNER_KEY":                      signerKey,
	}

	t.k8s.CreateSecret(namespace, "titan-secret", stringData)
}
