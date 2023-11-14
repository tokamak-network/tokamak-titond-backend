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
		p := getResourcesPath(tt.name)
		uts.Equal(p, uts.testDataPath+"/"+tt.name, "it is should be equal")
	}
}

func (uts *UtilsTestSuite) TestGetYAMLfiles() {
	tests := []struct {
		filePath  string
		resources []string
	}{
		{path.Join(uts.testDataPath, "l2geth"), []string{"Service", "StatefulSet"}},
	}

	for _, tt := range tests {
		for i, yamlfile := range getYAMLfiles(tt.filePath) {
			jsonBytes, _ := yaml.ToJSON(yamlfile)
			object, _ := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
			uObject, _ := object.(*unstructured.Unstructured)

			uts.Equal(tt.resources[i], uObject.GetKind())
		}
	}
}

func (uts *UtilsTestSuite) TestConvertYAMLtoObject() {
	tests := []struct {
		filePath  string
		resources []string
	}{
		{path.Join(uts.testDataPath, "l2geth"), []string{"Service", "StatefulSet"}},
	}

	for _, tt := range tests {
		for i, yamlfile := range getYAMLfiles(tt.filePath) {
			object := convertYAMLtoObject(yamlfile)
			uts.Equal(tt.resources[i], object.GetObjectKind().GroupVersionKind().Kind)
		}
	}
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
