package kubernetes

import (
	"context"
	"fmt"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *core.ConfigMap) (*core.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, v1.CreateOptions{})
}

func (k *Kubernetes) DeleteConfigMap(namespace string, name string) error {
	return k.client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
}

func (k *Kubernetes) GetConfigMap(namespace string, name string) (*core.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
}

func (k *Kubernetes) CreateConfigMapForApp(namespace string, rpc string, targetNetwork string, deployKey string) {
	fmt.Println("namespace", namespace, "rpc", rpc, "network", targetNetwork, "key", deployKey)
	configMapName := "deployer-configmap"
	_, err := k.GetConfigMap(namespace, configMapName)
	if err != nil {
		object, err := BuildObjectFromYamlFile("./deployments/deployer/configmap.yaml")
		if err != nil {
			panic(err)
		}
		configMap, _ := ConvertToConfigMap(object)
		UpdateConfigMapObjectName(configMap, configMapName)
		UpdateConfigMapObjectValue(configMap, "CONTRACTS_RPC_URL", rpc)
		UpdateConfigMapObjectValue(configMap, "CONTRACTS_TARGET_NETWORK", targetNetwork)
		UpdateConfigMapObjectValue(configMap, "CONTRACTS_DEPLOYER_KEY", deployKey)
		var configMapCreationErr error
		for i := 0; i < 5; i++ {
			_, configMapCreationErr = k.CreateConfigMap(namespace, configMap)
			if configMapCreationErr == nil {
				break
			}
		}
		if configMapCreationErr != nil {
			panic("Cannot init configMap for K8s cluster")
		}
	}
}

func (k *Kubernetes) CreateDeployerJob(namespace string) (*core.Pod, error) {
	return nil, nil
}

func (k *Kubernetes) DeleteDeployerJob(namespace string, name string) error {
	return nil
}
