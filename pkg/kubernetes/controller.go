package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateStatefulSet(namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return k.client.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateService(namespace string, service *corev1.Service) (*corev1.Service, error) {
	return k.client.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
}

func (k *Kubernetes) CreatePersistentVolumeClaim(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return k.client.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateIngress(namespace string, ingress *networkv1.Ingress) (*networkv1.Ingress, error) {
	return k.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
}
