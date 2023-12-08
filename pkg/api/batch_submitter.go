package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateBatchSubmitter(bs *model.Component) (*model.Component, error) {
	network, err := t.db.ReadNetwork(bs.NetworkID)
	if err != nil {
		return nil, err
	}

	l2geth, err := t.db.ReadComponentByType("l2geth", network.ID)
	if err != nil {
		return nil, err
	}
	if err := checkDependency(l2geth.Status); err != nil {
		return nil, err
	}

	result, err := t.db.CreateComponent(bs)
	if err == nil {
		go t.createBatchSubmitter(result, network.ContractAddressURL, t.config.ContractsRpcUrl)
	}

	return result, err
}

func (t *TitondAPI) createBatchSubmitter(bs *model.Component, contractAddressURL, l1RPC string) {
	namespace := generateNamespace(bs.NetworkID)
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "batch-submitter", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createBatchSubmitter error: convertToConfigmap")
		return
	}

	batchSubmitterConfig := map[string]string{
		"URL":        contractAddressURL,
		"L1_ETH_RPC": l1RPC,
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, batchSubmitterConfig)
	if err != nil {
		fmt.Printf("createBatchSubmitter error: %s\n", err)
		return
	}
	fmt.Printf("Created Batch-Submitter ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "batch-submitter", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		fmt.Printf("createBatchSubmitter error: convertToService")
		return
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		fmt.Printf("createBatchSubmitter error: %s\n", err)
		return
	}
	fmt.Printf("Created Batch-Submitter Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "batch-submitter", "deployment")
	deployment, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		fmt.Printf("createBatchSubmitter error: convertToDeployment")
		return
	}

	createdDeployment, err := t.k8s.CreateDeployment(namespace, deployment)
	if err != nil {
		fmt.Printf("createBatchSubmitter error: %s\n", err)
		return
	}
	fmt.Printf("Created Batch-Submitter Deployment: %s\n", createdDeployment.GetName())

	err = t.k8s.WaitingDeploymentCreated(createdDeployment.Namespace, createdDeployment.Name)
	if err != nil {
		/*TODO : rollback ? */
		return
	}
	bs.Status = true
	_, err = t.db.UpdateComponent(bs)
	if err != nil {
		return
	}

	return
}
