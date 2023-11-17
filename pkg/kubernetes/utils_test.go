package kubernetes

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type UtilsTestSuite struct {
	suite.Suite
	testDataPath string
}

func (uts *UtilsTestSuite) SetupSuite() {
	cPath, _ := os.Getwd()
	uts.testDataPath = path.Join(cPath, "../../testdata")
}

func (uts *UtilsTestSuite) TestGetResourcesPath() {
	tests := []struct {
		name string
	}{
		{"deployer"},
		{"l2geth"},
	}
	for _, tt := range tests {
		p := getDirPath(tt.name)
		uts.Equal(uts.testDataPath+"/"+tt.name, p, "it is should be equal")
	}
}

func (uts *UtilsTestSuite) TestGetYAMLfile() {
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
		uts.Equal(tt.resource, uObject.GetKind())
	}
}

func (uts *UtilsTestSuite) TestConvertYAMLtoObject() {
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
		uts.Equal(tt.resource, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
