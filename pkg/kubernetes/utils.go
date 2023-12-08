package kubernetes

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime/debug"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

func getDirPath(basePath, dirName string) string {
	rPath := path.Join(basePath, dirName)

	return rPath
}

func GetYAMLfile(basePath, dirName, fileName string) []byte {
	filePath := getDirPath(basePath, dirName)

	data, err := os.ReadFile(path.Join(filePath, fileName+".yaml"))
	if err != nil {
		path, _ := os.Getwd()
		fmt.Println("File path: ", filePath, "| Path:", path)
		debug.PrintStack()
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

func GetObject(manifestPath, componentName, yamlFileName string) runtime.Object {
	f := GetYAMLfile(manifestPath, componentName, yamlFileName)
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

func ConvertToIngress(obj runtime.Object) (*networkv1.Ingress, bool) {
	ingress, ok := obj.(*networkv1.Ingress)
	return ingress, ok
}

func BuildObjectFromYamlFile(file string) (runtime.Object, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sch := runtime.NewScheme()
	if err := scheme.AddToScheme(sch); err != nil {
		fmt.Println(err)
		return nil, err
	}

	decoder := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
	obj, _, err := decoder(yamlFile, nil, nil)
	return obj, err
}

func ConvertToConfigMap(obj runtime.Object) (*corev1.ConfigMap, bool) {
	instance, exist := obj.(*corev1.ConfigMap)
	return instance, exist
}

func ConvertToDeployment(obj runtime.Object) (*appsv1.Deployment, bool) {
	instance, exist := obj.(*appsv1.Deployment)
	return instance, exist
}

func UpdateConfigMapObjectValue(configMap *corev1.ConfigMap, key string, value string) {
	configMap.Data[key] = value
}
