package kubernetes

import (
	"context"
	"fmt"
	"time"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *core.ConfigMap) (*core.ConfigMap, error) {
	fmt.Println("Create configmap:", configMap)
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, v1.CreateOptions{})
}

func (k *Kubernetes) GetConfigMap(namespace string, name string) (*core.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
}

func (k *Kubernetes) CreateDeployment(namespace string, deployment *apps.Deployment) (*apps.Deployment, error) {
	fmt.Println("Create deployment:", deployment)
	return k.client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, v1.CreateOptions{})
}

func (k *Kubernetes) DeleteDeployment(namespace string, name string) error {
	fmt.Println("Delete deployment:", name)
	return k.client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
}

func (k *Kubernetes) CreateNamespace(name string) (*core.Namespace, error) {
	fmt.Println("Create namespace: ", name)
	namespace := &core.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	}
	return k.client.CoreV1().Namespaces().Create(context.TODO(), namespace, v1.CreateOptions{})
}

func (k *Kubernetes) GetNamespace(name string) (*core.Namespace, error) {
	return k.client.CoreV1().Namespaces().Get(context.TODO(), name, v1.GetOptions{})
}

func (k *Kubernetes) CreateNamespaceForApp(name string) {

	_, err := k.GetNamespace(name)
	if err != nil {
		for i := 0; i < 5; i++ {
			_, err := k.CreateNamespace(name)
			if err == nil {
				break
			}
		}
	}
}

func (k *Kubernetes) CreateConfigMapForApp(namespace string, rpc string, targetNetwork string, deployKey string) {
	configMapName := "deployer"
	fmt.Println("ConfigMap override data")
	fmt.Println("CONTRACTS_RPC_URL", rpc, len(rpc))
	fmt.Println("CONTRACTS_TARGET_NETWORK", targetNetwork, len(targetNetwork))
	fmt.Println("CONTRACTS_DEPLOYER_KEY", deployKey, len(deployKey))
	_, err := k.GetConfigMap(namespace, configMapName)
	if err != nil {
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

func (k *Kubernetes) GetPodsOfDeployment(namespace string, deployment string) (*core.PodList, error) {
	pods, err := k.client.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment),
	})
	return pods, err
}

func (k *Kubernetes) WaitingDeploymentCreated(namespace string, name string) error {
	var err error
	for i := 0; i < 60; i++ {
		deploy, err := k.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if deploy.Status.AvailableReplicas == 1 {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
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
	return nil
}
