package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func TestCreateComponent(t *testing.T) {
	// TODO: add more cases
	testcases := []struct {
		networkID     uint
		componentType string
		mockDBErr     error
		expectedErr   error
	}{
		// Case 1: create l2geth failed with db err
		{
			networkID:     1,
			componentType: "l2geth",
			mockDBErr:     errors.New("l2geth faied"),
			expectedErr:   errors.New("l2geth faied"),
		},
		// Case 2: create data-transport-layer failed with db err
		{
			networkID:     1,
			componentType: "data-transport-layer",
			mockDBErr:     errors.New("data-transport-layer faied"),
			expectedErr:   errors.New("data-transport-layer faied"),
		},
		// Case 3: create batch-submitter failed with db err
		{
			networkID:     1,
			componentType: "batch-submitter",
			mockDBErr:     errors.New("batch-submitter faied"),
			expectedErr:   errors.New("batch-submitter faied"),
		},
		// Case 4: create failed with type err
		{
			networkID:     1,
			componentType: "wrong type",
			mockDBErr:     nil,
			expectedErr:   apptypes.ErrInvalidComponentType,
		},
	}
	k8sClient := &MockK8sClient{}
	dbClient := &MockDBClient{}
	titond := NewTitondAPI(k8sClient, dbClient, nil, &Config{})

	for _, testcase := range testcases {
		component := &model.Component{
			Type:      testcase.componentType,
			NetworkID: testcase.networkID,
		}
		dbClient.err = testcase.mockDBErr
		_, err := titond.CreateComponent(component)
		assert.Equal(t, testcase.expectedErr, err)
	}
}

func TestGetComponentByType(t *testing.T) {
	// TODO: This is a predefined testing function
	k8sClient := &MockK8sClient{}
	titond := NewTitondAPI(k8sClient, nil, nil, &Config{})
	_, err := titond.GetComponentByType(12, "l2geth")
	assert.Equal(t, nil, err)
}

func TestGetComponentByID(t *testing.T) {
	// TODO: This is a predefined testing function
	k8sClient := &MockK8sClient{}
	titond := NewTitondAPI(k8sClient, nil, nil, &Config{})
	_, err := titond.GetComponentById(12)
	assert.Equal(t, nil, err)
}

func TestDeleteComponentByID(t *testing.T) {
	// TODO: This is a predefined testing function
	k8sClient := &MockK8sClient{}
	titond := NewTitondAPI(k8sClient, nil, nil, &Config{})
	err := titond.DeleteComponentById(12)
	assert.Equal(t, nil, err)
}
