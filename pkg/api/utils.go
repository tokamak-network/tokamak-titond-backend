package api

import "fmt"

func generateMountPath(name string, networkID, componentID uint) string {
	return fmt.Sprintf("%s-%d-%d", name, networkID, componentID)
}

func generateNamespace(networkID uint) string {
	return fmt.Sprintf("namespace-%d", networkID)
}
