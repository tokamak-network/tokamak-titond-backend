package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreate(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	fakeKubernetes := &Kubernetes{fakeClient}

	t.Run("Create StatefulSet", func(t *testing.T) {
		t.Run("Not StatefulSet yaml file", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "service")
			obj := ConvertYAMLtoObject(file)

			_, err := fakeKubernetes.CreateStatefulSet("test", obj)

			assert.EqualError(t, err, "This is not StatefulSet")
		})

		t.Run("Create StatefulSet Successfully", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "statefulset")
			obj := ConvertYAMLtoObject(file)

			res, err := fakeKubernetes.CreateStatefulSet("test", obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create Service", func(t *testing.T) {
		t.Run("Not Service yaml file", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "statefulset")
			obj := ConvertYAMLtoObject(file)

			_, err := fakeKubernetes.CreateService("test", obj)

			assert.EqualError(t, err, "This is not Service")
		})

		t.Run("Create Service Successfully", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "service")
			obj := ConvertYAMLtoObject(file)

			res, err := fakeKubernetes.CreateService("test", obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-svc", res.GetName(), "must be l2geth-svc")
		})
	})

	t.Run("Create PersistentVolumeClaim", func(t *testing.T) {
		t.Run("Not PersistentVolumeClaim yaml file", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "statefulset")
			obj := ConvertYAMLtoObject(file)

			_, err := fakeKubernetes.CreatePersistentVolumeClaim("test", obj)

			assert.EqualError(t, err, "This is not PersistentVolumeClaim")
		})

		t.Run("Create PersistentVolumeClaim Successfully", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "pvc")
			obj := ConvertYAMLtoObject(file)

			res, err := fakeKubernetes.CreatePersistentVolumeClaim("test", obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-pvc", res.GetName(), "must be l2geth-pvc")
		})
	})

	t.Run("Create ConfilgMap", func(t *testing.T) {
		t.Run("Not ConfilgMap yaml file", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "statefulset")
			obj := ConvertYAMLtoObject(file)

			_, err := fakeKubernetes.CreateConfigMap("test", obj)

			assert.EqualError(t, err, "This is not ConfilgMap")
		})

		t.Run("Create ConfilgMap Successfully", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "configMap")
			obj := ConvertYAMLtoObject(file)

			res, err := fakeKubernetes.CreateConfigMap("test", obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create Ingress", func(t *testing.T) {
		t.Run("Not Ingress yaml file", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "statefulset")
			obj := ConvertYAMLtoObject(file)

			_, err := fakeKubernetes.CreateIngress("test", obj)

			assert.EqualError(t, err, "This is not Ingress")
		})

		t.Run("Create Ingress Successfully", func(t *testing.T) {
			file := GetYAMLfile("l2geth", "ingress")
			obj := ConvertYAMLtoObject(file)

			res, err := fakeKubernetes.CreateIngress("test", obj)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "ingress-l2geth", res.GetName(), "must be ingress-l2geth")
		})
	})
}
