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

	result, err := t.db.CreateComponent(dtl)
	if err != nil {
		return nil, err
	}

	go t.createDTL(result, network.ContractAddressURL, t.config.ContractsRpcUrl)

	return result, nil
}

func (t *TitondAPI) createDTL(dtl *model.Component, contractAddressURL, l1RPC string) {
	namespace := generateNamespace(dtl.NetworkID)
	volumePath := generateVolumePath("dtl", dtl.NetworkID, dtl.ID)

	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "data-transport-layer", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToConfigmap")
		return
	}

	dtlConfig := map[string]string{
		"URL":                                   contractAddressURL,
		"DATA_TRANSPORT_LAYER__L1_RPC_ENDPOINT": l1RPC,
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, dtlConfig)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "data-transport-layer", "pv")
	pv, ok := kubernetes.ConvertToPersistentVolume(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToPersistentVolume")
		return
	}
	createdPV, err := t.k8s.CreatePersistentVolume(namespace, "dtl", pv)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer PersistentVolume: %s\n", createdPV.GetName())

	obj = kubernetes.GetObject(mPath, "data-transport-layer", "pvc")
	pvc, ok := kubernetes.ConvertToPersistentVolumeClaim(obj)
	if !ok {
		fmt.Printf("createDTL error: convertToPersistentVolumeClaim")
		return
	}
	createdPVC, err := t.k8s.CreatePersistentVolumeClaimWithAppSelector(namespace, "dtl", pvc)
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

	sfs.Spec.Template.Spec.Containers[0].VolumeMounts[0].SubPath = volumePath

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		fmt.Printf("createDTL error: %s\n", err)
		return
	}
	fmt.Printf("Created Data-Transport-Layer StatefulSet: %s\n", createdSFS.GetName())

	err = t.k8s.WatingStatefulsetCreated(createdSFS.Namespace, createdSFS.Name)
	if err != nil {
		return
	}
	dtl.Status = true
	_, err = t.db.UpdateComponent(dtl)

	if err != nil {
		/* TODO: rollback ? */
		return
	}
}
