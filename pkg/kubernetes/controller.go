package kubernetes

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateStatefulSet(namespace string, object runtime.Object) (*appsv1.StatefulSet, error) {
	statefulSet, ok := object.(*appsv1.StatefulSet)
	if !ok {
		err := errors.New("This is not StatefulSet")
		return nil, err
	}

	return k.client.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateService(namespace string, object runtime.Object) (*corev1.Service, error) {
	service, ok := object.(*corev1.Service)
	if !ok {
		err := errors.New("This is not Service")
		return nil, err
	}

	return k.client.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
}

func (k *Kubernetes) CreatePersistentVolumeClaim(namespace string, object runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := object.(*corev1.PersistentVolumeClaim)
	if !ok {
		err := errors.New("This is not PersistentVolumeClaim")
		return nil, err
	}

	return k.client.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateConfigMap(namespace string, object runtime.Object) (*corev1.ConfigMap, error) {
	configMap, ok := object.(*corev1.ConfigMap)
	if !ok {
		err := errors.New("This is not ConfilgMap")
		return nil, err
	}

	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateIngress(namespace string, object runtime.Object) (*networkv1.Ingress, error) {
	ingress, ok := object.(*networkv1.Ingress)
	if !ok {
		err := errors.New("This is not Ingress")
		return nil, err
	}

	return k.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
}
