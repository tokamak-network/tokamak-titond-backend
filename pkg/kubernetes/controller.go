package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
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

func (k *Kubernetes) CreateConfigMapWithConfig(namespace string, configMap *corev1.ConfigMap, config map[string]string) (*corev1.ConfigMap, error) {
	if configMap.Data == nil {
		configMap.Data = map[string]string{}
	}
	for k, v := range config {
		configMap.Data[k] = v
	}
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateIngress(namespace string, ingress *networkv1.Ingress) (*networkv1.Ingress, error) {
	return k.client.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (k *Kubernetes) UpdateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
}

func (k *Kubernetes) GetConfigMap(namespace string, name string) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (k *Kubernetes) CreateDeploymentWithName(namespace string, deployment *appsv1.Deployment, name string) (*appsv1.Deployment, error) {
	deployment.Name = name
	deployment.Spec.Selector.MatchLabels = map[string]string{"app": name}
	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{"app": name}
	return k.CreateDeployment(namespace, deployment)
}

func (k *Kubernetes) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return k.client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
}

func (k *Kubernetes) GetDeployment(namespace string, name string) (*appsv1.Deployment, error) {
	return k.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (k *Kubernetes) DeleteDeployment(namespace string, name string) error {
	return k.client.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (k *Kubernetes) CreateSecret(namespace, name string, stringData map[string]string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	if secret.StringData == nil {
		secret.StringData = map[string]string{}
	}

	for k, v := range stringData {
		secret.StringData[k] = v
	}

	secret.SetName(name)

	return k.client.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

func (k *Kubernetes) GetSecret(namespace, name string) (*corev1.Secret, error) {
	return k.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (k *Kubernetes) CreateNamespace(name string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return k.client.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
}

func (k *Kubernetes) GetNamespace(name string) (*corev1.Namespace, error) {
	return k.client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

func (k *Kubernetes) GetFileFromPod(namespace string, pod *corev1.Pod, path string) (string, error) {
	var err error
	addressCmd := []string{"cat", path}
	for i := 0; i < 200; i++ {
		stdout, stderr, err := k.Exec(namespace, pod, addressCmd)
		if err == nil && len(stderr) == 0 {
			return string(stdout), nil
		} else {
			if err == nil {
				err = errors.New(string(stderr))
			}
		}
		fmt.Println("Retry...", err)
		time.Sleep(time.Second * 10)
	}
	return "", err
}

func (k *Kubernetes) CreateConfigmapWithConfig(namespace string, template string, items map[string]string) error {
	object, err := BuildObjectFromYamlFile(template)
	if err != nil {
		return err
	}
	configMap, success := ConvertToConfigMap(object)
	if !success {
		return fmt.Errorf("cann't convert %s to a configmap", template)
	}
	for key, value := range items {
		configMap.Data[key] = value
	}
	_, err = k.GetConfigMap(namespace, configMap.Name)
	exist := (err != nil)
	var configMapCreationErr error
	for i := 0; i < 5; i++ {
		if exist {
			_, configMapCreationErr = k.CreateConfigMap(namespace, configMap)
			if configMapCreationErr == nil {
				break
			}
		} else {
			_, configMapCreationErr = k.UpdateConfigMap(namespace, configMap)
			if configMapCreationErr == nil {
				break
			}
		}
		time.Sleep(time.Second * 3)
	}
	return configMapCreationErr
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
	pods, err := k.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deployment),
	})
	return pods, err
}

func (k *Kubernetes) WaitingDeploymentCreated(namespace string, name string) error {
	var err error
	for i := 0; i < 1800; i++ {
		deploy, err := k.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
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

func (k *Kubernetes) WatingStatefulsetCreated(namespace, name string) error {
	var err error
	for i := 0; i < 1800; i++ {
		sfs, err := k.client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if sfs.Status.AvailableReplicas == 1 {
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
