package api

import (
	"errors"
	"fmt"
	"reflect"
)

func MakeDeployerName(id uint) string {
	return fmt.Sprintf("deployer-%d", id)
}

func generateMountPath(name string, networkID, componentID uint) string {
	return fmt.Sprintf("%s-%d-%d", name, networkID, componentID)
}

func generateNamespace(networkID uint) string {
	return fmt.Sprintf("namespace-%d", networkID)
}

func ConvertStructToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Struct {
		return nil, errors.New("the type of object is not struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()
		result[fieldName] = fieldValue
	}
	return result, nil
}
