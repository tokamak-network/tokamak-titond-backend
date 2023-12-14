package api

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestInternalCreateRelayer(t *testing.T) {
	testcases := []struct {
		relayer           *model.Component
		l1RPC             string
		addressFileURL    string
		mockConfigMap     *corev1.ConfigMap
		mockConfigMapErr  error
		mockDeployment    *appsv1.Deployment
		mockDeploymentErr error
		mockService       *corev1.Service
		mockServiceErr    error
		mockDBErr         error
		expectedExistErr  bool
	}{
		// Case 1: Happy case
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:             "l1RPC",
			addressFileURL:    "addressFileURL",
			mockConfigMap:     &corev1.ConfigMap{},
			mockConfigMapErr:  nil,
			mockDeployment:    &appsv1.Deployment{},
			mockDeploymentErr: nil,
			mockService:       &corev1.Service{},
			mockServiceErr:    nil,
			mockDBErr:         nil,
			expectedExistErr:  false,
		},
		// Case 2: Failed at making configmap on K8s
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:             "l1RPC",
			addressFileURL:    "addressFileURL",
			mockConfigMap:     &corev1.ConfigMap{},
			mockConfigMapErr:  errors.New("mock configmap err"),
			mockDeployment:    &appsv1.Deployment{},
			mockDeploymentErr: nil,
			mockService:       &corev1.Service{},
			mockServiceErr:    nil,
			mockDBErr:         nil,
			expectedExistErr:  true,
		},
		// Case 3: Failed at making deployment on K8s
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:             "l1RPC",
			addressFileURL:    "addressFileURL",
			mockConfigMap:     &corev1.ConfigMap{},
			mockConfigMapErr:  nil,
			mockDeployment:    &appsv1.Deployment{},
			mockDeploymentErr: errors.New("mock deployment err"),
			mockService:       &corev1.Service{},
			mockServiceErr:    nil,
			mockDBErr:         nil,
			expectedExistErr:  true,
		},
		// Case 4: Failed at making service on K8s
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:             "l1RPC",
			addressFileURL:    "addressFileURL",
			mockConfigMap:     &corev1.ConfigMap{},
			mockConfigMapErr:  nil,
			mockDeployment:    &appsv1.Deployment{},
			mockDeploymentErr: nil,
			mockService:       &corev1.Service{},
			mockServiceErr:    errors.New("mock service err"),
			mockDBErr:         nil,
			expectedExistErr:  true,
		},
		// Case 5: Failed at inserting to DB
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:             "l1RPC",
			addressFileURL:    "addressFileURL",
			mockConfigMap:     &corev1.ConfigMap{},
			mockConfigMapErr:  nil,
			mockDeployment:    &appsv1.Deployment{},
			mockDeploymentErr: nil,
			mockService:       &corev1.Service{},
			mockServiceErr:    nil,

			mockDBErr:        errors.New("mock db err"),
			expectedExistErr: true,
		},
	}
	path, _ := os.Getwd()
	path = path[0 : len(path)-len("/pkg/api")]
	os.Chdir(path)

	k8sClient := &MockK8sClient{
		manifestPath:   testDataPath,
		fileContent:    make(map[string]string),
		fileContentErr: make(map[string]error),
	}
	dbClient := &MockDBClient{}

	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		k8sClient.configmap = testcase.mockConfigMap
		k8sClient.configmapErr = testcase.mockConfigMapErr
		k8sClient.deployment = testcase.mockDeployment
		k8sClient.deploymentErr = testcase.mockDeploymentErr
		k8sClient.service = testcase.mockService
		k8sClient.serviceErr = testcase.mockServiceErr
		dbClient.componentUpdateErr = testcase.mockDBErr
		err := titond.createRelayer(testcase.relayer, testcase.l1RPC, testcase.addressFileURL)
		assert.Equal(t, testcase.expectedExistErr, err != nil)
	}
}

func TestCreateRelayer(t *testing.T) {
	testcases := []struct {
		relayer                *model.Component
		l1RPC                  string
		addressFileURL         string
		mockNetwork            *model.Network
		mockNetworkErr         error
		mockComponent          *model.Component
		mockComponentErr       error
		mockConfigMap          *corev1.ConfigMap
		mockConfigMapErr       error
		mockDeployment         *appsv1.Deployment
		mockDeploymentErr      error
		mockService            *corev1.Service
		mockServiceErr         error
		mockComponentUpdateErr error
		expectedExistErr       bool
	}{
		// Case 1: Happy case
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:          "l1RPC",
			addressFileURL: "addressFileURL",
			mockNetwork: &model.Network{
				ID: 12,
			},
			mockNetworkErr: nil,
			mockComponent: &model.Component{
				Type:      "l2geth",
				Status:    true,
				NetworkID: 12,
			},
			mockComponentErr:       nil,
			mockConfigMap:          &corev1.ConfigMap{},
			mockConfigMapErr:       nil,
			mockDeployment:         &appsv1.Deployment{},
			mockDeploymentErr:      nil,
			mockService:            &corev1.Service{},
			mockServiceErr:         nil,
			mockComponentUpdateErr: nil,
			expectedExistErr:       false,
		},
		// Case 2: Failed at query network
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:          "l1RPC",
			addressFileURL: "addressFileURL",
			mockNetwork: &model.Network{
				ID: 12,
			},
			mockNetworkErr: errors.New("query network failed"),
			mockComponent: &model.Component{
				Type:      "l2geth",
				Status:    true,
				NetworkID: 12,
			},
			mockComponentErr:       nil,
			mockConfigMap:          &corev1.ConfigMap{},
			mockConfigMapErr:       nil,
			mockDeployment:         &appsv1.Deployment{},
			mockDeploymentErr:      nil,
			mockService:            &corev1.Service{},
			mockServiceErr:         nil,
			mockComponentUpdateErr: nil,
			expectedExistErr:       true,
		},
		// Case 3: Failed at query l2geth
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:          "l1RPC",
			addressFileURL: "addressFileURL",
			mockNetwork: &model.Network{
				ID: 12,
			},
			mockNetworkErr: nil,
			mockComponent: &model.Component{
				Type:      "l2geth",
				Status:    true,
				NetworkID: 12,
			},
			mockComponentErr:       errors.New("query l2geth failed"),
			mockConfigMap:          &corev1.ConfigMap{},
			mockConfigMapErr:       nil,
			mockDeployment:         &appsv1.Deployment{},
			mockDeploymentErr:      nil,
			mockService:            &corev1.Service{},
			mockServiceErr:         nil,
			mockComponentUpdateErr: nil,
			expectedExistErr:       true,
		},
		// Case 4: Failed at query l2geth's status
		{
			relayer: &model.Component{
				Name:      "Relayer Component",
				Type:      "relayer",
				NetworkID: 12,
			},
			l1RPC:          "l1RPC",
			addressFileURL: "addressFileURL",
			mockNetwork: &model.Network{
				ID: 12,
			},
			mockNetworkErr: nil,
			mockComponent: &model.Component{
				Type:      "l2geth",
				Status:    false,
				NetworkID: 12,
			},
			mockComponentErr:       nil,
			mockConfigMap:          &corev1.ConfigMap{},
			mockConfigMapErr:       nil,
			mockDeployment:         &appsv1.Deployment{},
			mockDeploymentErr:      nil,
			mockService:            &corev1.Service{},
			mockServiceErr:         nil,
			mockComponentUpdateErr: nil,
			expectedExistErr:       true,
		},
	}
	path, _ := os.Getwd()
	path = path[0 : len(path)-len("/pkg/api")]
	os.Chdir(path)

	k8sClient := &MockK8sClient{
		manifestPath:   testDataPath,
		fileContent:    make(map[string]string),
		fileContentErr: make(map[string]error),
	}
	dbClient := &MockDBClient{}

	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		k8sClient.configmap = testcase.mockConfigMap
		k8sClient.configmapErr = testcase.mockConfigMapErr
		k8sClient.deployment = testcase.mockDeployment
		k8sClient.deploymentErr = testcase.mockDeploymentErr
		k8sClient.service = testcase.mockService
		k8sClient.serviceErr = testcase.mockServiceErr
		dbClient.network = testcase.mockNetwork
		dbClient.networkErr = testcase.mockNetworkErr
		dbClient.component = testcase.mockComponent
		dbClient.componentErr = testcase.mockComponentErr
		dbClient.componentUpdateErr = testcase.mockComponentUpdateErr
		_, err := titond.CreateRelayer(testcase.relayer)
		assert.Equal(t, testcase.expectedExistErr, err != nil)
	}
}
