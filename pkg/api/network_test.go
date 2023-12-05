package api

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testDataPath = "./testdata"

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
		manifestPath:   testDataPath,
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
		dbClient.networkErr = testcase.mockDBErr
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
		dbClient.networkErr = testcase.mockErrorFromDB
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
	dbClient.networkErr = nil
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
		dbClient.networkErr = testcase.mockErrorFromDB
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
		manifestPath:   testDataPath,
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
		k8sClient := &MockK8sClient{manifestPath: testDataPath}
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
	assert.Equal(t, nil, titond.updateDBWithValue(&model.Network{}, "", "", nil, nil))
}
