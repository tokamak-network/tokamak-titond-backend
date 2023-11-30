package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateL2Geth(l2geth *model.Component, config *ComponentConfig) (*model.Component, error) {
	network, err := t.db.ReadNetwork(l2geth.NetworkID)
	if err != nil {
		return nil, err
	}

	config.Namespace = generateNamespace(l2geth.NetworkID)
	config.Data["ROLLUP_STATE_DUMP_PATH"] = network.StateDumpURL

	result, err := t.db.CreateComponent(l2geth)
	if err == nil {
		go t.createL2Geth(config)
	}

	return result, err
}

/*
TODO :
  - config PV mount path
  - config public url
*/
func (t *TitondAPI) createL2Geth(config *ComponentConfig) {
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "l2geth", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToConfigmap")
		return
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(config.Namespace, cm, config.Data)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "l2geth", "pvc")
	pvc, ok := kubernetes.ConvertToPersistentVolumeClaim(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToPersistentVolumeClaim")
		return
	}

	createdPVC, err := t.k8s.CreatePersistentVolumeClaim(config.Namespace, pvc)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth PersistentVolumeClaim: %s\n", createdPVC.GetName())

	obj = kubernetes.GetObject(mPath, "l2geth", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToService")
		return
	}

	createdSVC, err := t.k8s.CreateService(config.Namespace, svc)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "l2geth", "statefulset")
	sfs, ok := kubernetes.ConvertToStatefulSet(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToStatefulSet")
		return
	}

	createdSFS, err := t.k8s.CreateStatefulSet(config.Namespace, sfs)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth StatefulSet: %s\n", createdSFS.GetName())

	obj = kubernetes.GetObject(mPath, "l2geth", "ingress")
	ingress, ok := kubernetes.ConvertToIngress(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToIngress")
		return
	}

	createdIngress, err := t.k8s.CreateIngress(config.Namespace, ingress)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth Ingress: %s\n", createdIngress.GetName())
}
