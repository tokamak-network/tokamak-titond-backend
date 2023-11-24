package kubernetes

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestGetDirPath(t *testing.T) {
	cPath, _ := os.Getwd()
	testDataPath := path.Join(cPath, "../../testdata")

	tests := []struct {
		name string
	}{
		{"deployer"},
		{"l2geth"},
	}

	for _, tt := range tests {
		p := getDirPath(tt.name)
		assert.Equal(t, testDataPath+"/"+tt.name, p, "it is should be equal")
	}
}

func TestGetYAMLfile(t *testing.T) {
	tests := []struct {
		dirName  string
		fileName string
		kind     string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		yamlfile := GetYAMLfile(tt.dirName, tt.fileName)
		jsonBytes, _ := yaml.ToJSON(yamlfile)
		object, _ := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
		uObject, _ := object.(*unstructured.Unstructured)
		assert.Equal(t, tt.kind, uObject.GetKind())
	}
}

func TestConvertBytestoObject(t *testing.T) {
	tests := []struct {
		dirName  string
		fileName string
		kind     string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		b := GetYAMLfile(tt.dirName, tt.fileName)
		obj := ConvertBytestoObject(b)
		assert.Equal(t, tt.kind, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}

func TestGetObject(t *testing.T) {
	tests := []struct {
		componentName string
		yamlFileName  string
		kind          string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		obj := GetObject(tt.componentName, tt.yamlFileName)
		assert.Equal(t, tt.kind, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}

func TestConvertToSpecificObject(t *testing.T) {
	t.Run("Convert to StatefulSet", func(t *testing.T) {
		t.Run("Not StatefulSet yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "service")
			_, ok := ConvertToStatefulSet(obj)

			assert.False(t, ok)
		})

		t.Run("Create StatefulSet Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")

			res, ok := ConvertToStatefulSet(obj)

			assert.True(t, ok)
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Convert to Service", func(t *testing.T) {
		t.Run("Not Service yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")
			_, ok := ConvertToService(obj)

			assert.False(t, ok)
		})

		t.Run("Convert Service Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "service")

			res, ok := ConvertToService(obj)

			assert.True(t, ok)
			assert.Equal(t, "l2geth-svc", res.GetName(), "must be l2geth-svc")
		})
	})

	t.Run("Convert to ConfigMap", func(t *testing.T) {
		t.Run("Not ConfigMap yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")
			_, ok := ConvertToConfigMap(obj)

			assert.False(t, ok)
		})

		t.Run("Convert ConfigMap Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "configMap")

			res, ok := ConvertToConfigMap(obj)

			assert.True(t, ok)
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Convert to PersistentVolumeClaim", func(t *testing.T) {
		t.Run("Not PersistentVolumeClaim yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "service")
			_, ok := ConvertToPersistentVolumeClaim(obj)

			assert.False(t, ok)
		})

		t.Run("Convert PersistentVolumeClaim Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "pvc")

			res, ok := ConvertToPersistentVolumeClaim(obj)

			assert.True(t, ok)
			assert.Equal(t, "l2geth-pvc", res.GetName(), "must be l2geth-pvc")
		})
	})

	t.Run("Convert to Ingress", func(t *testing.T) {
		t.Run("Not Ingress yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "pvc")
			_, ok := ConvertToIngress(obj)

			assert.False(t, ok)
		})

		t.Run("Convert Ingress Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "ingress")

			res, ok := ConvertToIngress(obj)

			assert.True(t, ok)
			assert.Equal(t, "ingress-l2geth", res.GetName(), "must be ingress-l2geth")
		})
	})
}

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
