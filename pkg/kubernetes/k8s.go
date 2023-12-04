package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	InCluster      bool
	KubeconfigPath string
	ManifestPath   string
}

type IK8s interface {
	GetPodStatus(namespace, name string) (string, error)
	CreateStatefulSet(namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error)
	CreateService(namespace string, service *corev1.Service) (*corev1.Service, error)
	CreatePersistentVolumeClaim(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error)
	CreateConfigMapWithConfig(namespace string, configMap *corev1.ConfigMap, config map[string]string) (*corev1.ConfigMap, error)
	CreateIngress(namespace string, ingress *networkv1.Ingress) (*networkv1.Ingress, error)
	CreateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error)
	GetConfigMap(namespace string, name string) (*corev1.ConfigMap, error)
	UpdateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error)
	CreateDeploymentWithName(namespace string, deployment *appsv1.Deployment, name string) (*appsv1.Deployment, error)
	CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	DeleteDeployment(namespace string, name string) error
	CreateNamespace(name string) (*corev1.Namespace, error)
	GetNamespace(name string) (*corev1.Namespace, error)
	GetFileFromPod(namespace string, pod *corev1.Pod, path string) (string, error)
	CreateConfigmapWithConfig(namespace string, template string, items map[string]string) error
	CreateNamespaceForApp(name string)
	GetPodsOfDeployment(namespace string, deployment string) (*corev1.PodList, error)
	WaitingDeploymentCreated(namespace string, name string) error
	Exec(namespace string, pod *corev1.Pod, command []string) ([]byte, []byte, error)
}

type Kubernetes struct {
	client       kubernetes.Interface
	config       *rest.Config
	manifestPath string
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

	return &Kubernetes{client, config, cfg.ManifestPath}, nil
}

func (k *Kubernetes) GetManifestPath() string {
	return k.manifestPath
}

func NewFakeKubernetes() *Kubernetes {
	fakeClient := fake.NewSimpleClientset()
	return &Kubernetes{fakeClient, nil, "../../testdata"}
}
