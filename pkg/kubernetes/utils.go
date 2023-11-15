package kubernetes

import (
	"io/ioutil"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

func BuildObjectFromYamlFile(file string) (runtime.Object, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sch := runtime.NewScheme()
	if err := scheme.AddToScheme(sch); err != nil {
		return nil, err
	}

	decoder := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, _, err := decoder(yamlFile, nil, nil)
	return obj, err
}

func ConvertToConfigMap(obj runtime.Object) (*core.ConfigMap, error) {
	return obj.(*core.ConfigMap), nil
}

func UpdateConfigMapObjectValue(configMap *core.ConfigMap, key string, value string) {
	configMap.Data[key] = value
}

func UpdateConfigMapObjectName(configMap *core.ConfigMap, value string) {
	configMap.Name = value
}
