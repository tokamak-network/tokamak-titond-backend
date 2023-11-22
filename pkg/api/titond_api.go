package api

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
)

type Config struct {
	Namespace              string
	ContractsRpcUrl        string
	ContractsTargetNetwork string
	ContractsDeployerKey   string
}

type TitondAPI struct {
	k8s         *kubernetes.Kubernetes
	db          db.Client
	fileManager services.IFIleManager
	config      *Config
}

func NewTitondAPI(k8s *kubernetes.Kubernetes, db db.Client, fileManager services.IFIleManager, config *Config) *TitondAPI {
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
		panic("Cannot init configmap for deployer")
	}
}
