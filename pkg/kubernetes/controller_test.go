package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
)

func TestCreate(t *testing.T) {
	fakeKubernetes := NewFakeKubernetes()
	mPath := fakeKubernetes.GetManifestPath()

	t.Run("Create StatefulSet", func(t *testing.T) {
		t.Run("Create StatefulSet Successfully", func(t *testing.T) {
			obj := GetObject(mPath, "l2geth", "statefulset")
			sfs, _ := ConvertToStatefulSet(obj)

			res, err := fakeKubernetes.CreateStatefulSet("test", sfs)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create Service", func(t *testing.T) {
		t.Run("Create Service Successfully", func(t *testing.T) {
			obj := GetObject(mPath, "l2geth", "service")
			svc, _ := ConvertToService(obj)

			res, err := fakeKubernetes.CreateService("test", svc)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-svc", res.GetName(), "must be l2geth-svc")
		})
	})

	t.Run("Create PersistentVolumeClaim", func(t *testing.T) {
		t.Run("Create PersistentVolumeClaim Successfully", func(t *testing.T) {
			obj := GetObject(mPath, "l2geth", "pvc")
			pvc, _ := ConvertToPersistentVolumeClaim(obj)

			res, err := fakeKubernetes.CreatePersistentVolumeClaim("test", pvc)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-pvc", res.GetName(), "must be l2geth-pvc")
		})
	})

	t.Run("Create ConfilgMap", func(t *testing.T) {
		t.Run("Create ConfilgMap Successfully", func(t *testing.T) {
			obj := GetObject(mPath, "l2geth", "configMap")
			configMap, _ := ConvertToConfigMap(obj)

			res, err := fakeKubernetes.CreateConfigMap("test", configMap)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth", res.GetName(), "must be l2geth")
		})
	})

	t.Run("Create ConfigMap With Config", func(t *testing.T) {
		obj := GetObject(mPath, "l2geth", "configMap")
		tCm, _ := ConvertToConfigMap(obj)

		tests := []struct {
			name     string
			cm       *corev1.ConfigMap
			data     map[string]string
			expected map[string]string
		}{
			{
				name:     "Test empty config",
				cm:       &corev1.ConfigMap{},
				data:     map[string]string{},
				expected: map[string]string{},
			},
			{
				name: "Test exist config",
				cm:   &corev1.ConfigMap{},
				data: map[string]string{
					"testKey": "testValue",
				},
				expected: map[string]string{
					"testKey": "testValue",
				},
			},
			{
				name: "Test change data",
				cm:   tCm,
				data: map[string]string{
					"ETH1_CONFIRMATION_DEPTH": "10",
					"GASPRICE":                "100",
				},
				expected: map[string]string{
					"ETH1_CONFIRMATION_DEPTH": "10",
					"GASPRICE":                "100",
				},
			},
		}

		for i, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.cm.SetName(fmt.Sprintf("TestConfigMap_%d", i))
				cm, _ := fakeKubernetes.CreateConfigMapWithConfig("default", tt.cm, tt.data)
				assert.Equal(t, fmt.Sprintf("TestConfigMap_%d", i), cm.GetName())
				if len(tt.cm.Data) > 0 {
					for k := range tt.expected {
						assert.Equal(t, tt.expected[k], cm.Data[k])
					}
				}
			})
		}
	})

	t.Run("Create Ingress", func(t *testing.T) {
		t.Run("Create Ingress Successfully", func(t *testing.T) {
			obj := GetObject(mPath, "l2geth", "ingress")
			ingress, _ := ConvertToIngress(obj)

			res, err := fakeKubernetes.CreateIngress("test", ingress)

			assert.NoError(t, err, "must be not error")
			assert.Equal(t, "l2geth-ingress", res.GetName(), "must be l2geth-ingress")
		})
	})
}
