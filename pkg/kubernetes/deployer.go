package kubernetes

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func (k *Kubernetes) GetDeployerResult(namespace string, pod *core.Pod) (string, error, string, error) {
	fmt.Println("Get Deploy result")
	addresses, addressError := k.GetFileFromPod(namespace, pod, "/opt/optimism/packages/tokamak/contracts/genesis/addresses.json")
	stateDump, stateError := k.GetFileFromPod(namespace, pod, "/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json")
	return addresses, addressError, stateDump, stateError
}

func (k *Kubernetes) CreateDeployer(namespace string, name string) (*apps.Deployment, error) {
	fmt.Println("Create deployer: ", name)
	object, _ := BuildObjectFromYamlFile("./deployments/deployer/deployment.yaml")
	deployment, success := ConvertToDeployment(object)
	if !success {
		panic("Failed to convert to deployment object")
	}
	deployment.Name = name
	deployment.Spec.Selector.MatchLabels = map[string]string{"app": name}
	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{"app": name}

	var deployerCreationErr error
	for i := 0; i < 5; i++ {
		_, deployerCreationErr = k.CreateDeployment(namespace, deployment)
		if deployerCreationErr == nil {
			deployerCreationErr = k.WaitingDeploymentCreated(namespace, name)
			if deployerCreationErr == nil {
				break
			}
		}
	}
	return deployment, deployerCreationErr
}

func (k *Kubernetes) DeleteDeployer(namespace string, name string) error {
	return k.DeleteDeployment(namespace, name)
}

func (k *Kubernetes) CreateConfigMapForDeployer(namespace string, rpc string, targetNetwork string, deployKey string) {
	overrideData := map[string]string{
		"CONTRACTS_RPC_URL":        rpc,
		"CONTRACTS_TARGET_NETWORK": targetNetwork,
		"CONTRACTS_DEPLOYER_KEY":   deployKey,
	}
	err := k.OverrideConfigMapFromTemplate(namespace, "./deployments/deployer/configmap.yaml", overrideData)
	if err != nil {
		panic("Cannot init configmap for deployer")
	}
}
