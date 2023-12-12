package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateRelayer(relayer *model.Component) (*model.Component, error) {
	network, err := t.db.ReadNetwork(relayer.NetworkID)
	if err != nil {
		return nil, err
	}

	l2geth, err := t.db.ReadComponentByType("l2geth", relayer.NetworkID)
	if err != nil {
		return nil, err
	}

	if err := checkDependency(l2geth.Status); err != nil {
		return nil, err
	}

	result, err := t.db.CreateComponent(relayer)
	if err == nil {
		go t.createRelayer(result, t.config.ContractsRpcUrl, network.ContractAddressURL)
	}

	return result, err

}

func (t *TitondAPI) createRelayer(relayer *model.Component, l1RPC string, addressFileUrl string) error {
	namespace := generateNamespace(relayer.NetworkID)
	publicURL := generatePublcURL("l2geth", relayer.NetworkID, relayer.ID)

	obj, _ := kubernetes.BuildObjectFromYamlFile("./deployments/relayer/configmap.yaml")
	configMapObj, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		panic("createRelayer error: convertToConfigmap")
	}
	relayerConfig := map[string]string{
		"URL":                              addressFileUrl,
		"MESSAGE_RELAYER__L1_RPC_PROVIDER": l1RPC,
	}

	configMap, err := t.k8s.CreateConfigMapWithConfig(namespace, configMapObj, relayerConfig)
	if err != nil {
		fmt.Println("createRelayer error:", err)
		return err
	}
	fmt.Println("Created Relayer ConfigMap:", configMap.GetName())

	obj, _ = kubernetes.BuildObjectFromYamlFile("./deployments/relayer/service.yaml")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		panic("createRelayer error: convertToService")
	}

	service, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		fmt.Println("createRelayer error:", err)
		return err
	}
	fmt.Println("Created Relayer Service:", service.GetName())

	obj, _ = kubernetes.BuildObjectFromYamlFile("./deployments/relayer/deployment.yaml")
	deploymentObj, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		panic("createRelayer error: convertToDeployment")
	}

	deployment, err := t.k8s.CreateDeployment(namespace, deploymentObj)
	if err != nil {
		fmt.Printf("createRelayer error: %s\n", err)
		return err
	}
	fmt.Printf("Created Relayer Deployment: %s\n", deployment.GetName())

	err = t.k8s.WaitingDeploymentCreated(deployment.Namespace, deployment.Name)
	if err != nil {
		return err
	}
	relayer.Status = true

	obj, _ = kubernetes.BuildObjectFromYamlFile("./deployments/relayer/ingress.yaml")
	ingressObj, ok := kubernetes.ConvertToIngress(obj)
	if !ok {
		panic("createRelayer error: convertToIngress")
	}

	ingressObj.Spec.Rules[0].Host = publicURL

	ingress, err := t.k8s.CreateIngress(namespace, ingressObj)
	if err != nil {
		fmt.Println("createRelayer error:", err)
		return err
	}
	fmt.Println("Created Relayer Ingress:", ingress.GetName(), " ||| URL: ", publicURL)
	relayer.PublicURL = publicURL

	_, err = t.db.UpdateComponent(relayer)
	return err

}
