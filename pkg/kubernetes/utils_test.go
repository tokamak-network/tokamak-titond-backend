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
	testDataPath    string
	deploymentsPath string
}

func (uts *UtilsTestSuite) SetupTest() {
	cPath, _ := os.Getwd()
	uts.testDataPath = path.Join(cPath, "../../testdata")
	uts.deploymentsPath = path.Join(cPath, "../../deployments")
}

func (uts *UtilsTestSuite) TestGetResourcesPath() {
	tests := []struct {
		name string
	}{
		{"deployer"},
		{"l2geth"},
	}
	for _, tt := range tests {
		p := getResourcesPath(tt.name)
		uts.Equal(p, uts.deploymentsPath+"/"+tt.name, "it is should be equal")
	}
}

func (uts *UtilsTestSuite) TestGetYAMLfiles() {
	tests := []struct {
		name      string
		filePath  string
		resources []string
	}{
		{"l2geth", path.Join(uts.testDataPath, "l2geth"), []string{"Service", "StatefulSet"}},
	}

	for _, tt := range tests {
		for i, yamlfile := range getYAMLfiles(tt.name, tt.filePath) {
			jsonBytes, _ := yaml.ToJSON(yamlfile)
			object, _ := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
			uObject, _ := object.(*unstructured.Unstructured)

			uts.Equal(tt.resources[i], uObject.GetKind())
		}
	}
}

func (uts *UtilsTestSuite) TestConvertYAMLtoObject() {
	tests := []struct {
		name      string
		filePath  string
		resources []string
	}{
		{"l2geth", path.Join(uts.testDataPath, "l2geth"), []string{"Service", "StatefulSet"}},
	}

	for _, tt := range tests {
		for i, yamlfile := range getYAMLfiles(tt.name, tt.filePath) {
			object := convertYAMLtoObject(yamlfile)
			uts.Equal(tt.resources[i], object.GetObjectKind().GroupVersionKind().Kind)
		}
	}
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
