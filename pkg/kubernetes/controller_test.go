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
		t.Run("Create StatefulSet Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "statefulset")
			sfs, _ := ConvertToStatefulSet(obj)

			res, err := fakeKubernetes.CreateStatefulSet("test", sfs)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create Service", func(t *testing.T) {
		t.Run("Create Service Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "service")
			svc, _ := ConvertToService(obj)

			res, err := fakeKubernetes.CreateService("test", svc)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-svc", res.GetName(), "must be l2geth-svc")
		})
	})

	t.Run("Create PersistentVolumeClaim", func(t *testing.T) {
		t.Run("Create PersistentVolumeClaim Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "pvc")
			pvc, _ := ConvertToPersistentVolumeClaim(obj)

			res, err := fakeKubernetes.CreatePersistentVolumeClaim("test", pvc)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-pvc", res.GetName(), "must be l2geth-pvc")
		})
	})

	t.Run("Create ConfilgMap", func(t *testing.T) {
		t.Run("Create ConfilgMap Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "configMap")
			configMap, _ := ConvertToConfigMap(obj)

			res, err := fakeKubernetes.CreateConfigMap("test", configMap)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create Ingress", func(t *testing.T) {
		t.Run("Create Ingress Successfully", func(t *testing.T) {
			obj := GetObject("l2geth", "ingress")
			ingress, _ := ConvertToIngress(obj)

			res, err := fakeKubernetes.CreateIngress("test", ingress)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "ingress-l2geth", res.GetName(), "must be ingress-l2geth")
		})
	})
}
