package kubernetes

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateTitondNode(nodeName, namespace string) {
	r, _ := GetResourceFactory(nodeName)
	r.Create(k.client)
}

func (k *Kubernetes) CreateL2Geth(namespace, name string) error {
	rPath := getResourcesPath("l2geth")
	yamls := getYAMLfiles(rPath)

	for _, yaml := range yamls {
		object := convertYAMLtoObject(yaml).(*unstructured.Unstructured)
		kind := object.GetObjectKind().GroupVersionKind().Kind

		switch kind {
		case "Service":
			var svc corev1.Service
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &svc)
			if err != nil {
				return err
			}

			svc.SetNamespace(namespace)
			svc.SetGenerateName("l2geth")
			svc.SetName(name)

			created, err := k.client.CoreV1().Services(namespace).Create(context.TODO(), &svc, metav1.CreateOptions{})
			if err != nil {
				return err
			}

			log.Printf("L2Geth Service Created %s/%s\n", namespace, created.GetName())

		case "StatefulSet":
			var sts appsv1.StatefulSet
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &sts)
			if err != nil {
				return err
			}

			sts.SetNamespace(namespace)
			sts.SetGenerateName("l2geth")
			sts.SetName(name)

			created, err := k.client.AppsV1().StatefulSets(namespace).Create(context.TODO(), &sts, metav1.CreateOptions{})
			if err != nil {
				return err
			}

			log.Printf("L2Geth StatefulSet Created %s/%s\n", namespace, created.GetName())
		}

	}

	return nil
}
