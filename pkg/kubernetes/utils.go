package kubernetes

import (
	"fmt"
	"io/ioutil"

	app "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
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

func ConvertToConfigMap(obj runtime.Object) (*core.ConfigMap, bool) {
	instance, exist := obj.(*core.ConfigMap)
	return instance, exist
}

func ConvertToDeployment(obj runtime.Object) (*app.Deployment, bool) {
	instance, exist := obj.(*app.Deployment)
	return instance, exist
}

func UpdateConfigMapObjectValue(configMap *core.ConfigMap, key string, value string) {
	configMap.Data[key] = value
}
