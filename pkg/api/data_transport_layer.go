package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateDTL(dtl *model.Component) (*model.Component, error) {
	network, err := t.db.ReadNetwork(dtl.NetworkID)
	if err != nil {
		return nil, err
	}

	namespace := generateNamespace(dtl.NetworkID)
	contractAddressURL := network.ContractAddressURL

	result, err := t.db.CreateComponent(dtl)
	if err == nil {
		go t.createDTL(namespace, contractAddressURL)
	}

	return result, err
}

func (t *TitondAPI) createDTL(namespace, contractAddressURL string) {
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "data-transport-layer", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToConfigmap")
		return
	}

	dtlConfig := map[string]string{
		"URL": contractAddressURL,
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, dtlConfig)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "data-transport-layer", "pvc")
	pvc, ok := kubernetes.ConvertToPersistentVolumeClaim(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToPersistentVolumeClaim")
		return
	}

	createdPVC, err := t.k8s.CreatePersistentVolumeClaim(namespace, pvc)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer PersistentVolumeClaim: %s\n", createdPVC.GetName())

	obj = kubernetes.GetObject(mPath, "data-transport-layer", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToService")
		return
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "data-transport-layer", "statefulset")
	sfs, ok := kubernetes.ConvertToStatefulSet(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToStatefulSet")
		return
	}

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer StatefulSet: %s\n", createdSFS.GetName())
}
