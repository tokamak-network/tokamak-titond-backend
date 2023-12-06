package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeDeployerName(t *testing.T) {
	var id uint = 12
	result := MakeDeployerName(id)
	assert.Equal(t, "deployer-12", result)
}

func TestGenerateMountPath(t *testing.T) {
	name := "l2geth"
	networkID := 1
	componentID := 1
	result := generateMountPath(name, uint(networkID), uint(componentID))
	assert.Equal(t, "l2geth-1-1", result)
}

func TestGenerateNamespace(t *testing.T) {
	var id uint = 12
	result := generateNamespace(id)
	assert.Equal(t, "namespace-12", result)
}
