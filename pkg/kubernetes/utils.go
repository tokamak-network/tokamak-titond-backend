package kubernetes

import (
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

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
