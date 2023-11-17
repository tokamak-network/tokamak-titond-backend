package components

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type L2Geth struct {
	name      string
	namespace string

	l1RPC         string
	confirmations string

	resources map[string]runtime.Object
}

func NewL2Geth(options ...func(*L2Geth)) *L2Geth {
	l2geth := &L2Geth{
		name:      "default-names",
		namespace: "default",
	}

	l2geth.resources = map[string]runtime.Object{
		"statefulset": &appsv1.StatefulSet{},
		"service":     &corev1.Service{},
		"configMap":   &corev1.ConfigMap{},
	}

	for _, option := range options {
		option(l2geth)
	}

	return l2geth
}

func WithName(name string) func(*L2Geth) {
	return func(l2geth *L2Geth) {
		l2geth.name = name
	}
}

func WithNamespace(namespace string) func(*L2Geth) {
	return func(l2geth *L2Geth) {
		l2geth.namespace = namespace
	}
}

func WithFiles(resources ...string) func(*L2Geth) {
	return func(l2geth *L2Geth) {
		for _, r := range resources {
			f := utils.GetYAMLfile("l2geth", r+".yaml")
			obj := utils.ConvertYAMLtoObject(f)
			l2geth.resources[r] = obj
		}
	}
}
