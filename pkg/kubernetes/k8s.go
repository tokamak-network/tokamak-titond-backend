package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	InCluster      bool
	KubeconfigPath string
}

type Kubernetes struct {
	client kubernetes.Interface
}

func NewKubernetes(cfg *Config) (*Kubernetes, error) {
	var config *rest.Config
	var err error
	if cfg.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Kubernetes{client}, nil
}

func NewFakeKubernetes() *Kubernetes {
	fakeClient := fake.NewSimpleClientset()
	return &Kubernetes{fakeClient}
}
