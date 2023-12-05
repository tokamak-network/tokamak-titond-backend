package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
)

type MockDBClient struct {
	network      *model.Network
	networks     []model.Network
	numOfDeleted int64
	component    *model.Component
	err          error
	networkID    uint
	offset       int
	limit        int
}

func (client *MockDBClient) CreateNetwork(network *model.Network) (*model.Network, error) {
	return client.network, client.err
}

func (client *MockDBClient) ReadNetwork(networkID uint) (*model.Network, error) {
	client.networkID = networkID
	return client.network, client.err
}

func (client *MockDBClient) ReadNetworkByRange(offset int, limit int) ([]model.Network, error) {
	client.offset = offset
	client.limit = limit
	return client.networks[offset : offset+limit], client.err
}

func (client *MockDBClient) UpdateNetwork(network *model.Network) (*model.Network, error) {
	return client.network, client.err
}

func (client *MockDBClient) DeleteNetwork(networkID uint) (int64, error) {
	return client.numOfDeleted, client.err
}

func (client *MockDBClient) CreateComponent(component *model.Component) (*model.Component, error) {
	return client.component, client.err
}

func (client *MockDBClient) ReadComponent() {

}

func (client *MockDBClient) ReadAllComponent() {

}

func (client *MockDBClient) UpdateComponent() {

}

func (client *MockDBClient) DeleteComponent() {

}

type MockK8sClient struct {
	err          error
	manifestPath string
	podStatus    string
	fileContent  string
	namespace    *corev1.Namespace
	podList      *corev1.PodList
	service      *corev1.Service
	pvc          *corev1.PersistentVolumeClaim
	configmap    *corev1.ConfigMap
	ingress      *networkv1.Ingress
	deployment   *appsv1.Deployment
	statefulSet  *appsv1.StatefulSet
	stdout       []byte
	stderr       []byte
}

func (client *MockK8sClient) GetManifestPath() string {
	return client.manifestPath
}

func (client *MockK8sClient) GetPodStatus(namespace, name string) (string, error) {
	return client.podStatus, client.err
}

func (client *MockK8sClient) CreateStatefulSet(namespace string, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return client.statefulSet, client.err
}

func (client *MockK8sClient) CreateService(namespace string, service *corev1.Service) (*corev1.Service, error) {
	return client.service, client.err
}

func (client *MockK8sClient) CreatePersistentVolumeClaim(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	return client.pvc, client.err
}

func (client *MockK8sClient) CreateConfigMapWithConfig(namespace string, configMap *corev1.ConfigMap, config map[string]string) (*corev1.ConfigMap, error) {
	return client.configmap, client.err
}

func (client *MockK8sClient) CreateIngress(namespace string, ingress *networkv1.Ingress) (*networkv1.Ingress, error) {
	return client.ingress, client.err
}

func (client *MockK8sClient) CreateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return client.configmap, client.err
}

func (client *MockK8sClient) GetConfigMap(namespace string, name string) (*corev1.ConfigMap, error) {
	return client.configmap, client.err
}

func (client *MockK8sClient) UpdateConfigMap(namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	return client.configmap, client.err
}

func (client *MockK8sClient) CreateDeploymentWithName(namespace string, deployment *appsv1.Deployment, name string) (*appsv1.Deployment, error) {
	return client.deployment, client.err
}

func (client *MockK8sClient) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return client.deployment, client.err
}

func (client *MockK8sClient) DeleteDeployment(namespace string, name string) error {
	return client.err
}

func (client *MockK8sClient) CreateNamespace(name string) (*corev1.Namespace, error) {
	return client.namespace, client.err
}

func (client *MockK8sClient) GetNamespace(name string) (*corev1.Namespace, error) {
	return client.namespace, client.err
}

func (client *MockK8sClient) GetFileFromPod(namespace string, pod *corev1.Pod, path string) (string, error) {
	return client.fileContent, client.err
}

func (client *MockK8sClient) CreateConfigmapWithConfig(namespace string, template string, items map[string]string) error {
	return client.err
}

func (client *MockK8sClient) CreateNamespaceForApp(name string) {

}

func (client *MockK8sClient) GetPodsOfDeployment(namespace string, deployment string) (*corev1.PodList, error) {
	return client.podList, client.err
}

func (client *MockK8sClient) WaitingDeploymentCreated(namespace string, name string) error {
	return client.err
}

func (client *MockK8sClient) Exec(namespace string, pod *corev1.Pod, command []string) ([]byte, []byte, error) {
	return client.stdout, client.stderr, client.err
}

type MockFileManager struct {
	fileName string
	content  string
	err      error
}

func (fileManager *MockFileManager) UploadContent(fileName string, content string) (string, error) {
	fileManager.fileName = fileName
	fileManager.content = content
	return fileName, fileManager.err
}

func TestCleanK8sJob(t *testing.T) {
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	fileManager := &MockFileManager{}
	titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})
	k8sClient.err = nil
	assert.Equal(t, nil, titond.cleanK8sJob(&model.Network{}))
	// t.Etitond.cleanK8sJob()

}
