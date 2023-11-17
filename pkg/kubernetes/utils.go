package kubernetes

import (
	"log"
	"os"
	"path"

	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

var mDir string = "../../deployments"

func getDirPath(dirName string) string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rPath := path.Join(currentPath, mDir, dirName)

	return rPath
}

func GetYAMLfile(dirName, fileName string) []byte {
	filePath := getDirPath(dirName)

	data, err := os.ReadFile(path.Join(filePath, fileName+".yaml"))
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func ConvertBytestoObject(b []byte) runtime.Object {
	jsonBytes, err := yaml.ToJSON(b)
	if err != nil {
		log.Fatal(err)
	}

	serializer := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory,
		scheme.Scheme,
		scheme.Scheme,
		jsonserializer.SerializerOptions{
			Yaml:   true,
			Pretty: false,
			Strict: false,
		},
	)

	object, err := runtime.Decode(serializer, jsonBytes)
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
