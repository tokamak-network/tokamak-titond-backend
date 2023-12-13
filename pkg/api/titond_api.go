package api

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
)

type Config struct {
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
	return &TitondAPI{
		k8s,
		db,
		fileManager,
		config,
	}
}
