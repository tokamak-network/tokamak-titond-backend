package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

func (k *Kubernetes) GetPodStatus(namespace, name string) (string, error) {
	pod, err := k.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
	return string(pod.Status.Phase), err
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, v1.CreateOptions{})
}

func (k *Kubernetes) UpdateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, v1.UpdateOptions{})
}

func (k *Kubernetes) GetConfigMap(namespace string, name string) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
}

func (k *Kubernetes) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return k.client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, v1.CreateOptions{})
}

func (k *Kubernetes) DeleteDeployment(namespace string, name string) error {
	return k.client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
}

func (k *Kubernetes) CreateNamespace(name string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
	}
	return k.client.CoreV1().Namespaces().Create(context.TODO(), namespace, v1.CreateOptions{})
}

func (k *Kubernetes) GetNamespace(name string) (*corev1.Namespace, error) {
	return k.client.CoreV1().Namespaces().Get(context.TODO(), name, v1.GetOptions{})
}

func (k *Kubernetes) CreateNamespaceForApp(name string) {
	_, err := k.GetNamespace(name)
	if err != nil {
		for i := 0; i < 5; i++ {
			_, err := k.CreateNamespace(name)
			if err == nil {
				break
			}
			time.Sleep(time.Second * 10)
		}
	}
}

func (k *Kubernetes) GetPodsOfDeployment(namespace string, deployment string) (*corev1.PodList, error) {
	pods, err := k.client.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment),
	})
	return pods, err
}

func (k *Kubernetes) WaitingDeploymentCreated(namespace string, name string) error {
	var err error
	for i := 0; i < 1800; i++ {
		deploy, err := k.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if deploy.Status.AvailableReplicas == 1 {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

func (k *Kubernetes) Exec(namespace string, pod *corev1.Pod, command []string) ([]byte, []byte, error) {
	if len(pod.Spec.Containers) == 0 {
		return nil, nil, errors.New("no container in the pod")
	}
	req := k.client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, nil, errors.New("err when adding to scheme")
	}
	var stdout bytes.Buffer
	paramCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Command:   command,
		Container: pod.Spec.Containers[0].Name,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	},
		paramCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())
	if err != nil {
		return nil, nil, errors.New("err when creating executor")
	}
	var stderr bytes.Buffer
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error stream cmd: %v", err)
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}
