package kubernetes

import (
	"fmt"
	"time"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

func (k *Kubernetes) GetDeployerResult(namespace string, pod *core.Pod) (string, string) {
	fmt.Println("Get Deploy result")
	var addresses string
	var stateDump string
	addressCmd := []string{"cat", "/opt/optimism/packages/tokamak/contracts/genesis/addresses.json"}
	for i := 0; i < 200; i++ {
		stdout, stderr, err := k.Exec(namespace, pod, addressCmd)
		if err == nil && len(stderr) == 0 {
			addresses = string(stdout)
			break
		}
		fmt.Println("Retry...", err)
		time.Sleep(time.Second * 10)
	}
	stateDumpCmd := []string{"cat", "/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json"}
	for i := 0; i < 100; i++ {
		stdout, stderr, err := k.Exec(namespace, pod, stateDumpCmd)
		if err == nil && len(stderr) == 0 {
			stateDump = string(stdout)
			break
		}
		fmt.Println("Retry...", err)
		time.Sleep(time.Second * 3)
	}
	return addresses, stateDump
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
