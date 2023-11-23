package kubernetes

import (
	"log"
	"os"
	"path"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
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

func GetObject(componentName, yamlFileName string) runtime.Object {
	f := GetYAMLfile(componentName, yamlFileName)
	return ConvertBytestoObject(f)
}

func ConvertToStatefulSet(obj runtime.Object) (*appsv1.StatefulSet, bool) {
	sfs, ok := obj.(*appsv1.StatefulSet)
	return sfs, ok
}

func ConvertToService(obj runtime.Object) (*corev1.Service, bool) {
	svc, ok := obj.(*corev1.Service)
	return svc, ok
}

func ConvertToPersistentVolumeClaim(obj runtime.Object) (*corev1.PersistentVolumeClaim, bool) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	return pvc, ok
}

func ConvertToConfigMap(obj runtime.Object) (*corev1.ConfigMap, bool) {
	configMap, ok := obj.(*corev1.ConfigMap)
	return configMap, ok
}

func ConvertToIngress(obj runtime.Object) (*networkv1.Ingress, bool) {
	ingress, ok := obj.(*networkv1.Ingress)
	return ingress, ok
}

func init() {
	if os.Getenv("MODE") == "test" {
		mDir = "../../testdata"
	}
}
