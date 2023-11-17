package kubernetes

import (
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
		componentName string
		fileName      string
		resource      string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
	}

	for _, tt := range tests {
		yamlfile := GetYAMLfile(tt.componentName, tt.fileName)
		jsonBytes, _ := yaml.ToJSON(yamlfile)
		object, _ := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
		uObject, _ := object.(*unstructured.Unstructured)
		assert.Equal(t, tt.resource, uObject.GetKind())
	}
}

func TestConvertYAMLtoObject(t *testing.T) {
	tests := []struct {
		componentName string
		fileName      string
		resource      string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
	}

	for _, tt := range tests {
		yamlfile := GetYAMLfile(tt.componentName, tt.fileName)
		obj := ConvertYAMLtoObject(yamlfile)
		assert.Equal(t, tt.resource, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}
