package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func TestGetComponentByID(t *testing.T) {
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})
	dbClient.component = &model.Component{}
	dbClient.networkErr = nil
	network, err := titond.GetComponentById(1)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, network)
}

func TestGetComponentByType(t *testing.T) {
	testcases := []struct {
		mockDataFromDB  *model.Component
		mockErrorFromDB error
		expectedErr     error
	}{
		{
			mockDataFromDB:  nil,
			mockErrorFromDB: errors.New("err"),
			expectedErr:     errors.New("err"),
		},
		{
			mockDataFromDB:  &model.Component{},
			mockErrorFromDB: nil,
			expectedErr:     nil,
		},
	}
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		dbClient.component = testcase.mockDataFromDB
		dbClient.componentErr = testcase.mockErrorFromDB
		component, err := titond.GetComponentByType(1, "")
		assert.Equal(t, testcase.expectedErr, err)
		assert.Equal(t, testcase.expectedErr == nil, component != nil)

	}
}

func TestDeleteComponent(t *testing.T) {
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
		dbClient.componentErr = testcase.mockErrorFromDB
		assert.Equal(t, testcase.expectedErr, titond.DeleteComponentById(1))
	}
}
