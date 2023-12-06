package api

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return client.networks, client.err

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
	err                 error
	deployerCreationErr error
	deployerWaitingErr  error
	getPodListErr       error
	manifestPath        string
	podStatus           string
	fileContent         map[string]string
	fileContentErr      map[string]error
	namespace           *corev1.Namespace
	podList             *corev1.PodList
	service             *corev1.Service
	pvc                 *corev1.PersistentVolumeClaim
	configmap           *corev1.ConfigMap
	ingress             *networkv1.Ingress
	deployment          *appsv1.Deployment
	statefulSet         *appsv1.StatefulSet
	stdout              []byte
	stderr              []byte
}

// func (client *MockK8sClient) clearMapData() {

// }

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
	return client.deployment, client.deployerCreationErr
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
	return client.fileContent[path], client.fileContentErr[path]
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
	return client.deployerWaitingErr
}

func (client *MockK8sClient) Exec(namespace string, pod *corev1.Pod, command []string) ([]byte, []byte, error) {
	return client.stdout, client.stderr, client.err
}

type MockFileManager struct {
	fileName string
	content  string
	url      map[string]string
	err      error
}

func (fileManager *MockFileManager) UploadContent(fileName string, content string) (string, error) {
	fileManager.fileName = fileName
	fileManager.content = content
	return fileManager.url[fileName], fileManager.err
}

func TestCreateNetwork(t *testing.T) {
	testcases := []struct {
		networkID               uint
		mockDBErr               error
		mockAddressData         string
		mockAddressURL          string
		mockDumpData            string
		mockDumpUrl             string
		mockPodListErr          error
		mockUploadAddressErr    error
		mockUploadDumpErr       error
		mockAddressErr          error
		mockDumpEr              error
		mockDeployerCreationErr error
		mockDeployerWaitingErr  error
		updateDBErr             error
		podList                 *corev1.PodList
		expectedErr             error
	}{
		// Case 1: happy case
		{
			networkID:               1,
			mockDBErr:               nil,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "deployer",
							Namespace: "titond",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			networkID:               1,
			mockDBErr:               errors.New("mock create err in db"),
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList:                 &corev1.PodList{},
			expectedErr:             errors.New("mock create err in db"),
		},
	}
	path, _ := os.Getwd()
	path = path[0 : len(path)-8]
	os.Chdir(path)

	k8sClient := &MockK8sClient{
		fileContent:    make(map[string]string),
		fileContentErr: make(map[string]error),
	}
	dbClient := &MockDBClient{}
	fileManager := &MockFileManager{
		url: make(map[string]string),
	}
	titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})

	for _, testcase := range testcases {

		network := &model.Network{ID: testcase.networkID}
		dbClient.network = network
		dbClient.err = testcase.mockDBErr
		k8sClient.deployerCreationErr = testcase.mockDeployerCreationErr
		k8sClient.deployerWaitingErr = testcase.mockDeployerWaitingErr
		k8sClient.podList = testcase.podList
		k8sClient.getPodListErr = testcase.mockPodListErr
		k8sClient.fileContent["/opt/optimism/packages/tokamak/contracts/genesis/addresses.json"] = testcase.mockAddressData
		k8sClient.fileContent["/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json"] = testcase.mockDumpData
		k8sClient.fileContentErr["/opt/optimism/packages/tokamak/contracts/genesis/addresses.json"] = testcase.mockAddressErr
		k8sClient.fileContentErr["/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json"] = testcase.mockDumpEr
		addressFileName := fmt.Sprintf("address-%d.json", network.ID)
		fileManager.url[addressFileName] = testcase.mockAddressURL
		dumpFileName := fmt.Sprintf("state-dump-%d.json", network.ID)
		fileManager.url[dumpFileName] = testcase.mockDumpUrl
		_, err := titond.CreateNetwork(network)
		assert.Equal(t, testcase.expectedErr, err)
	}
}

func TestGetNetworkByPage(t *testing.T) {
	testcases := []struct {
		mockDataFromDB  []model.Network
		mockErrorFromDB error
		expectedErr     error
	}{
		{
			mockDataFromDB:  nil,
			mockErrorFromDB: nil,
			expectedErr:     types.ErrResourceNotFound,
		},
		{
			mockDataFromDB:  []model.Network{},
			mockErrorFromDB: nil,
			expectedErr:     types.ErrResourceNotFound,
		},
		{
			mockDataFromDB:  nil,
			mockErrorFromDB: errors.New("err"),
			expectedErr:     errors.New("err"),
		},
		{
			mockDataFromDB: []model.Network{
				{ID: 1},
			},
			mockErrorFromDB: nil,
			expectedErr:     nil,
		},
	}
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		dbClient.networks = testcase.mockDataFromDB
		dbClient.err = testcase.mockErrorFromDB
		networks, err := titond.GetNetworksByPage(1)
		assert.Equal(t, testcase.expectedErr, err)
		if testcase.expectedErr == nil {
			assert.Greater(t, len(networks), 0)
		}
	}
}

func TestGetNetworkByID(t *testing.T) {
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})
	dbClient.network = &model.Network{}
	dbClient.err = nil
	network, err := titond.GetNetworkByID(1)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, network)
}

func TestDeleteNetwork(t *testing.T) {
	testcases := []struct {
		mockDataFromDB  int64
		mockErrorFromDB error
		expectedErr     error
	}{
		{
			mockDataFromDB:  0,
			mockErrorFromDB: nil,
			expectedErr:     types.ErrResourceNotFound,
		},
		{
			mockDataFromDB:  1,
			mockErrorFromDB: nil,
			expectedErr:     nil,
		},
		{
			mockDataFromDB:  0,
			mockErrorFromDB: errors.New("delete err"),
			expectedErr:     errors.New("delete err"),
		},
	}
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		dbClient.numOfDeleted = testcase.mockDataFromDB
		dbClient.err = testcase.mockErrorFromDB
		assert.Equal(t, testcase.expectedErr, titond.DeleteNetwork(1))
	}
}

func TestInternalCreateNetwork(t *testing.T) {
	testcases := []struct {
		networkID               uint
		mockAddressData         string
		mockAddressURL          string
		mockDumpData            string
		mockDumpUrl             string
		mockPodListErr          error
		mockUploadAddressErr    error
		mockUploadDumpErr       error
		mockAddressErr          error
		mockDumpEr              error
		mockDeployerCreationErr error
		mockDeployerWaitingErr  error
		updateDBErr             error
		podList                 *corev1.PodList
		expectedAddressURL      string
		expectedDumpURL         string
	}{
		// Case 1: happy case
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "deployer",
							Namespace: "titond",
						},
					},
				},
			},
			expectedAddressURL: "mockAddressURL",
			expectedDumpURL:    "mockDumpUrl",
		},
		// Case 2: no pod list returned
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{},
			},
			expectedAddressURL: "",
			expectedDumpURL:    "",
		},
		// Case 3: failed at deployer cannot create k8s's job
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: errors.New("mockDeployerCreationErr"),
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{},
			},
			expectedAddressURL: "",
			expectedDumpURL:    "",
		},
		// Case 4: failed at deployer cannot waiting result of k8s's job
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  errors.New("mockDeployerWaitingErr"),
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{},
			},
			expectedAddressURL: "",
			expectedDumpURL:    "",
		},
		// Case 5: no addressURL
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          errors.New("no address url"),
			mockDumpEr:              nil,
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "deployer",
							Namespace: "titond",
						},
					},
				},
			},
			expectedAddressURL: "",
			expectedDumpURL:    "mockDumpUrl",
		},
		// Case 6: no dumpURL
		{
			networkID:               1,
			mockAddressData:         "mockAddressData",
			mockAddressURL:          "mockAddressURL",
			mockDumpData:            "mockDumpData",
			mockDumpUrl:             "mockDumpUrl",
			mockPodListErr:          nil,
			mockUploadAddressErr:    nil,
			mockUploadDumpErr:       nil,
			mockAddressErr:          nil,
			mockDumpEr:              errors.New("no dump url"),
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			updateDBErr:             nil,
			podList: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "deployer",
							Namespace: "titond",
						},
					},
				},
			},
			expectedAddressURL: "mockAddressURL",
			expectedDumpURL:    "",
		},
	}

	path, _ := os.Getwd()
	path = path[0 : len(path)-8]
	os.Chdir(path)

	k8sClient := &MockK8sClient{
		fileContent:    make(map[string]string),
		fileContentErr: make(map[string]error),
	}
	dbClient := &MockDBClient{}
	fileManager := &MockFileManager{
		url: make(map[string]string),
	}
	titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})

	for _, testcase := range testcases {

		network := &model.Network{ID: testcase.networkID}
		k8sClient.deployerCreationErr = testcase.mockDeployerCreationErr
		k8sClient.deployerWaitingErr = testcase.mockDeployerWaitingErr
		k8sClient.podList = testcase.podList
		k8sClient.getPodListErr = testcase.mockPodListErr
		k8sClient.fileContent["/opt/optimism/packages/tokamak/contracts/genesis/addresses.json"] = testcase.mockAddressData
		k8sClient.fileContent["/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json"] = testcase.mockDumpData
		k8sClient.fileContentErr["/opt/optimism/packages/tokamak/contracts/genesis/addresses.json"] = testcase.mockAddressErr
		k8sClient.fileContentErr["/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json"] = testcase.mockDumpEr
		addressFileName := fmt.Sprintf("address-%d.json", network.ID)
		fileManager.url[addressFileName] = testcase.mockAddressURL
		dumpFileName := fmt.Sprintf("state-dump-%d.json", network.ID)
		fileManager.url[dumpFileName] = testcase.mockDumpUrl
		addressURL, dumpURL := titond.createNetwork(network)
		// t.Log("addressURL: ", addressURL, " dumpURL:", dumpURL)
		assert.Equal(t, testcase.expectedAddressURL, addressURL)
		assert.Equal(t, testcase.expectedDumpURL, dumpURL)
	}
}

func TestCreateDeployer(t *testing.T) {
	testcases := []struct {
		mockDeployerCreationErr error
		mockDeployerWaitingErr  error
		expectedErrExist        bool
	}{
		{
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  nil,
			expectedErrExist:        false,
		},
		{
			mockDeployerCreationErr: errors.New("mockDeployerCreationErr"),
			mockDeployerWaitingErr:  nil,
			expectedErrExist:        true,
		},
		{
			mockDeployerCreationErr: nil,
			mockDeployerWaitingErr:  errors.New("mockDeployerWaitingErr"),
			expectedErrExist:        true,
		},
	}

	path, _ := os.Getwd()
	path = path[0 : len(path)-8]
	os.Chdir(path)

	for _, testcase := range testcases {
		k8sClient := &MockK8sClient{}
		dbClient := &MockDBClient{}
		fileManager := &MockFileManager{}
		titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})

		k8sClient.deployerCreationErr = testcase.mockDeployerCreationErr
		k8sClient.deployerWaitingErr = testcase.mockDeployerWaitingErr
		_, err := titond.createDeployer("titond", "deployer")
		t.Log("Err: ", err)
		assert.Equal(t, testcase.expectedErrExist, err != nil)
	}
}

func TestUploadAddressFile(t *testing.T) {
	k8sClient := &MockK8sClient{}
	fileManager := &MockFileManager{
		url: make(map[string]string),
	}
	titond := NewTitondAPI(k8sClient, nil, fileManager, &Config{})
	fileManager.err = nil
	fileManager.url["address file"] = "url"
	url, err := titond.uploadAddressFile("address file", "address data")
	assert.Equal(t, "url", url)
	assert.Equal(t, nil, err)
}

func TestUploadDumpFile(t *testing.T) {
	k8sClient := &MockK8sClient{}
	fileManager := &MockFileManager{
		url: make(map[string]string),
	}
	titond := NewTitondAPI(k8sClient, nil, fileManager, &Config{})
	fileManager.err = nil
	fileManager.url["dump file"] = "url"
	url, err := titond.uploadDumpFile("dump file", "dump data")
	assert.Equal(t, "url", url)
	assert.Equal(t, nil, err)
}

func TestUpdateDBWithValue(t *testing.T) {
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	fileManager := &MockFileManager{}
	titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})
	k8sClient.err = nil
	assert.Equal(t, nil, titond.cleanK8sJob(&model.Network{}))
}

func TestCleanK8sJob(t *testing.T) {
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	fileManager := &MockFileManager{}
	titond := NewTitondAPI(k8sClient, dbClient, fileManager, &Config{})
	k8sClient.err = nil
	assert.Equal(t, nil, titond.updateDBWithValue(&model.Network{}, "", "", nil, nil))
}
