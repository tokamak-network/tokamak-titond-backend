package kubernetes

import (
	"errors"
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

func ConvertToStatefulSet(obj runtime.Object) (*appsv1.StatefulSet, error) {
	sfs, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		err := errors.New("this is not StatefulSet")
		return nil, err
	}

	return sfs, nil
}

func ConvertToService(obj runtime.Object) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		err := errors.New("this is not Service")
		return nil, err
	}

	return svc, nil
}

func ConvertToPersistentVolumeClaim(obj runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		err := errors.New("this is not PersistentVolumeClaim")
		return nil, err
	}

	return pvc, nil
}

func ConvertToConfigMap(obj runtime.Object) (*corev1.ConfigMap, error) {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		err := errors.New("this is not ConfigMap")
		return nil, err
	}

	return configMap, nil
}

func ConvertToIngress(obj runtime.Object) (*networkv1.Ingress, error) {
	ingress, ok := obj.(*networkv1.Ingress)
	if !ok {
		err := errors.New("this is not Ingress")
		return nil, err
	}

	return ingress, nil
}

func init() {
	if os.Getenv("MODE") == "test" {
		mDir = "../../testdata"
	}
}
