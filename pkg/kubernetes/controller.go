package kubernetes

import (
	"context"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) GetNamespace() (*core.NamespaceList, error) {
	namespaceList, err := k.client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	return namespaceList, err
}
