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
	deployment, _ := ConvertToDeployment(object)
	deployment.Name = name
	deployment.Spec.Selector.MatchLabels = map[string]string{"app": name}
	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{"app": name}

	var deployerCreationErr error
	for i := 0; i < 5; i++ {
		_, deployerCreationErr = k.CreateDeployment(namespace, deployment)
		if deployerCreationErr == nil {
			break
		}
	}
	return deployment, deployerCreationErr
}

func (k *Kubernetes) DeleteDeployer(namespace string, name string) error {
	return k.DeleteDeployment(namespace, name)
}

func (k *Kubernetes) CreateConfigMapForDeployer(namespace string, rpc string, targetNetwork string, deployKey string) {
	configMapName := "deployer"
	fmt.Println("ConfigMap override data")
	fmt.Println("CONTRACTS_RPC_URL", rpc, len(rpc))
	fmt.Println("CONTRACTS_TARGET_NETWORK", targetNetwork, len(targetNetwork))
	fmt.Println("CONTRACTS_DEPLOYER_KEY", deployKey, len(deployKey))
	_, err := k.GetConfigMap(namespace, configMapName)
	exist := (err != nil)
	object, err := BuildObjectFromYamlFile("./deployments/deployer/configmap.yaml")
	if err != nil {
		panic(err)
	}
	configMap, _ := ConvertToConfigMap(object)
	fmt.Println("Original configmap data")
	fmt.Println(" CONTRACTS_RPC_URL: ", configMap.Data["CONTRACTS_RPC_URL"], len(configMap.Data["CONTRACTS_RPC_URL"]))
	fmt.Println(" CONTRACTS_TARGET_NETWORK: ", configMap.Data["CONTRACTS_TARGET_NETWORK"], len(configMap.Data["CONTRACTS_TARGET_NETWORK"]))
	fmt.Println(" CONTRACTS_DEPLOYER_KEY: ", configMap.Data["CONTRACTS_DEPLOYER_KEY"], len(configMap.Data["CONTRACTS_DEPLOYER_KEY"]))
	UpdateConfigMapObjectValue(configMap, "CONTRACTS_RPC_URL", rpc)
	UpdateConfigMapObjectValue(configMap, "CONTRACTS_TARGET_NETWORK", targetNetwork)
	UpdateConfigMapObjectValue(configMap, "CONTRACTS_DEPLOYER_KEY", deployKey)
	fmt.Println("After override configmap data")
	fmt.Println(" CONTRACTS_RPC_URL: ", configMap.Data["CONTRACTS_RPC_URL"], len(configMap.Data["CONTRACTS_RPC_URL"]))
	fmt.Println(" CONTRACTS_TARGET_NETWORK: ", configMap.Data["CONTRACTS_TARGET_NETWORK"], len(configMap.Data["CONTRACTS_TARGET_NETWORK"]))
	fmt.Println(" CONTRACTS_DEPLOYER_KEY: ", configMap.Data["CONTRACTS_DEPLOYER_KEY"], len(configMap.Data["CONTRACTS_DEPLOYER_KEY"]))

	var configMapCreationErr error
	for i := 0; i < 5; i++ {
		if exist {
			_, configMapCreationErr = k.CreateConfigMap(namespace, configMap)
			if configMapCreationErr == nil {
				break
			}
		} else {
			_, configMapCreationErr = k.UpdateConfigMap(namespace, configMap)
			if configMapCreationErr == nil {
				break
			}
		}

	}
	if configMapCreationErr != nil {
		panic("Cannot init configMap for K8s cluster")
	}
}
