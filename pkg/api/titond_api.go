package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
)

type Config struct {
	Namespace              string
	ContractsRpcUrl        string
	ContractsTargetNetwork string
	ContractsDeployerKey   string
	SequencerKey           string
	ProposerKey            string
	RelayerWallet          string
	SignerKey              string
}

type ITitondAPI interface {
	CreateNetwork(data *model.Network) (*model.Network, error)
	GetNetworksByPage(page int) ([]model.Network, error)
	GetNetworkByID(networkID uint) (*model.Network, error)
	DeleteNetwork(id uint) error
	CreateComponent(component *model.Component) (*model.Component, error)
	GetComponentByType(networkID uint, componentType string) (*model.Component, error)
	GetComponentById(componentID uint) (*model.Component, error)
	DeleteComponentById(componentID uint) error
}

type TitondAPI struct {
	k8s         kubernetes.IK8s
	db          db.Client
	fileManager services.IFIleManager
	config      *Config
}

func NewTitondAPI(k8s kubernetes.IK8s, db db.Client, fileManager services.IFIleManager, config *Config) *TitondAPI {
	titondAPI := &TitondAPI{
		k8s,
		db,
		fileManager,
		config,
	}
	titondAPI.Initialize()
	return titondAPI
}

func (t *TitondAPI) Initialize() {
	t.k8s.CreateNamespaceForApp(t.config.Namespace)
	overrideData := map[string]string{
		"CONTRACTS_RPC_URL":        t.config.ContractsRpcUrl,
		"CONTRACTS_TARGET_NETWORK": t.config.ContractsTargetNetwork,
		"CONTRACTS_DEPLOYER_KEY":   t.config.ContractsDeployerKey,
	}
	err := t.k8s.CreateConfigmapWithConfig(t.config.Namespace, "./deployments/deployer/configmap.yaml", overrideData)
	if err != nil {
		fmt.Println("err: ", err)
		panic("Cannot init configmap for deployer")
	}
}
