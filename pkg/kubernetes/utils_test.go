package kubernetes

import (
	"fmt"
	"testing"
)

func TestBuildAPIObjectFromYamlFile(t *testing.T) {
	_, err := BuildObjectFromYamlFile("../../deployments/deployer/configmap.yaml")
	if err != nil {
		t.Error("Failed when decoding a configmap.yaml")
	}
	_, err = BuildObjectFromYamlFile("../../deployments/deployer/deployment.yaml")
	if err != nil {
		t.Error("Failed when decoding a deployment.yaml")
	}
}

func TestUpdateConfigMapObjectValue(t *testing.T) {
	object, err := BuildObjectFromYamlFile("../../deployments/deployer/configmap.yaml")
	fmt.Println(" [] ", object, err)
	if err != nil {
		t.Error("Failed when decoding a configmap.yaml")
	}
	configMap, exist := ConvertToConfigMap(object)
	if !exist {
		t.Error("Failed when converting to a configmap.yaml")
	}

	UpdateConfigMapObjectValue(configMap, "CONTRACTS_DEPLOYER_KEY", "0123456789")
	value, exist := configMap.Data["CONTRACTS_DEPLOYER_KEY"]
	if !exist {
		t.Error("CONTRACTS_DEPLOYER_KEY need to be exist")
	}
	if value != "0123456789" {
		t.Error("Update ConfigMap Value failed")
	}
}

func TestConvertToDeployment(t *testing.T) {
	object, err := BuildObjectFromYamlFile("../../deployments/deployer/deployment.yaml")
	if err != nil {
		t.Error("Failed when decoding a deployment.yaml")
	}
	_, exist := ConvertToConfigMap(object)
	if exist {
		t.Error("A deployment is converted to a configmap")
	}
	_, exist = ConvertToDeployment(object)
	if !exist {
		t.Error("Failed to convert to a deployment object")
	}
}
