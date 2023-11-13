package kubernetes

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) ApplyResources() {

}
