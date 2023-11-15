package kubernetes

import (
	"log"
	"os"
	"path"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var mDir string = "../../deployments"

func getResourcesPath(name string) string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rPath := path.Join(currentPath, mDir, name)

	return rPath
}

func getYAMLfiles(filePath string) [][]byte {
	var yamlFiles [][]byte

	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := os.ReadFile(path.Join(filePath, "/", file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		yamlFiles = append(yamlFiles, data)
	}

	return yamlFiles
}

func convertYAMLtoObject(yamlfile []byte) runtime.Object {
	jsonBytes, err := yaml.ToJSON(yamlfile)
	if err != nil {
		log.Fatal(err)
	}

	object, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	return object
}

func init() {
	if os.Getenv("MODE") == "test" {
		mDir = "../../testdata"
	}
}