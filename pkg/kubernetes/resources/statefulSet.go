package resources

import (
	"context"
	"errors"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type StatefulSet struct{}

func (sfs *StatefulSet) Create(client kubernetes.Interface, obj runtime.Object, setObject func(*appsv1.StatefulSet)) error {
	statefulset, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		err := errors.New("This is not statefulset")
		return err
	}

	setObject(statefulset)

	namespace := statefulset.GetNamespace()

	created, err := client.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulset, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Statefulset Created %s/%s\n", namespace, created.GetName())

	return nil
}
