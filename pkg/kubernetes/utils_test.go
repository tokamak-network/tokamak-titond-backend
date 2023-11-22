package kubernetes

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestGetDirPath(t *testing.T) {
	cPath, _ := os.Getwd()
	testDataPath := path.Join(cPath, "../../testdata")

	tests := []struct {
		name string
	}{
		{"deployer"},
		{"l2geth"},
	}

	for _, tt := range tests {
		p := getDirPath(tt.name)
		assert.Equal(t, testDataPath+"/"+tt.name, p, "it is should be equal")
	}
}

func TestGetYAMLfile(t *testing.T) {
	tests := []struct {
		dirName  string
		fileName string
		kind     string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		yamlfile := GetYAMLfile(tt.dirName, tt.fileName)
		jsonBytes, _ := yaml.ToJSON(yamlfile)
		object, _ := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
		uObject, _ := object.(*unstructured.Unstructured)
		assert.Equal(t, tt.kind, uObject.GetKind())
	}
}

func TestConvertBytestoObject(t *testing.T) {
	tests := []struct {
		dirName  string
		fileName string
		kind     string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		b := GetYAMLfile(tt.dirName, tt.fileName)
		obj := ConvertBytestoObject(b)
		assert.Equal(t, tt.kind, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}

func TestGetObject(t *testing.T) {
	tests := []struct {
		componentName string
		yamlFileName  string
		kind          string
	}{
		{"l2geth", "statefulset", "StatefulSet"},
		{"l2geth", "service", "Service"},
		{"l2geth", "configMap", "ConfigMap"},
		{"l2geth", "pvc", "PersistentVolumeClaim"},
		{"l2geth", "ingress", "Ingress"},
	}

	for _, tt := range tests {
		obj := GetObject(tt.componentName, tt.yamlFileName)
		assert.Equal(t, tt.kind, obj.GetObjectKind().GroupVersionKind().Kind)
	}
}

func TestConvertToSpecificObject(t *testing.T) {
	t.Run("Convert to StatefulSet", func(t *testing.T) {
		t.Run("Not StatefulSet yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "service")
			_, err := ConvertToStatefulSet(obj)

			assert.EqualError(t, err, "this is not StatefulSet")
		})

		t.Run("Create StatefulSet Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")

			res, err := ConvertToStatefulSet(obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Convert to Service", func(t *testing.T) {
		t.Run("Not Service yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")
			_, err := ConvertToService(obj)

			assert.EqualError(t, err, "this is not Service")
		})

		t.Run("Convert Service Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "service")

			res, err := ConvertToService(obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-svc", res.GetName(), "must be l2geth-svc")
		})
	})

	t.Run("Convert to ConfigMap", func(t *testing.T) {
		t.Run("Not ConfigMap yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")
			_, err := ConvertToConfigMap(obj)

			assert.EqualError(t, err, "this is not ConfigMap")
		})

		t.Run("Convert ConfigMap Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "configMap")

			res, err := ConvertToConfigMap(obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Convert to PersistentVolumeClaim", func(t *testing.T) {
		t.Run("Not PersistentVolumeClaim yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "service")
			_, err := ConvertToPersistentVolumeClaim(obj)

			assert.EqualError(t, err, "this is not PersistentVolumeClaim")
		})

		t.Run("Convert PersistentVolumeClaim Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "pvc")

			res, err := ConvertToPersistentVolumeClaim(obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-pvc", res.GetName(), "must be l2geth-pvc")
		})
	})

	t.Run("Convert to Ingress", func(t *testing.T) {
		t.Run("Not Ingress yaml file", func(t *testing.T) {
			obj := GetObject("l2geth", "pvc")
			_, err := ConvertToIngress(obj)

			assert.EqualError(t, err, "this is not Ingress")
		})

		t.Run("Convert Ingress Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "ingress")

			res, err := ConvertToIngress(obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "ingress-l2geth", res.GetName(), "must be ingress-l2geth")
		})
	})
}
