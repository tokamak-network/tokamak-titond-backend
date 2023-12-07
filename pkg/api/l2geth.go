package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateL2Geth(l2geth *model.Component) (*model.Component, error) {
	network, err := t.db.ReadNetwork(l2geth.NetworkID)
	if err != nil {
		return nil, err
	}
	if err := checkDependency(network.Status); err != nil {
		return nil, err
	}
	dtl, err := t.db.ReadComponentByType("data-transport-layer", network.ID)
	if err != nil {
		return nil, err
	}
	if err := checkDependency(dtl.Status); err != nil {
		return nil, err
	}

	result, err := t.db.CreateComponent(l2geth)
	if err != nil {
		return nil, err
	}
	go t.createL2Geth(result, network.StateDumpURL, t.config.ContractsRpcUrl)

	return result, nil
}

func (t *TitondAPI) createL2Geth(l2geth *model.Component, stateDumpURL, l1RPC string) {
	namespace := generateNamespace(l2geth.NetworkID)
	volumePath := generateVolumePath("l2geth", l2geth.NetworkID, l2geth.ID)
	publicURL := generatePublcURL("l2geth", l2geth.NetworkID, l2geth.ID)

	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "l2geth", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToConfigmap")
		return
	}

	l2gethConfig := map[string]string{
		"ROLLUP_STATE_DUMP_PATH": stateDumpURL,
		"ETH1_HTTP":              "",
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, l2gethConfig)
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

	createdPVC, err := t.k8s.CreatePersistentVolumeClaim(namespace, pvc)
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

	createdSVC, err := t.k8s.CreateService(namespace, svc)
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

	sfs.Spec.Template.Spec.Containers[0].VolumeMounts[0].SubPath = volumePath

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth StatefulSet: %s\n", createdSFS.GetName())

	err = t.k8s.WatingStatefulsetCreated(createdSFS.Namespace, createdSFS.Name)
	if err != nil {
		/*TODO : rollback ? */
		return
	}
	l2geth.Status = true

	obj = kubernetes.GetObject(mPath, "l2geth", "ingress")
	ingress, ok := kubernetes.ConvertToIngress(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToIngress")
		return
	}

	ingress.Spec.Rules[0].Host = publicURL

	createdIngress, err := t.k8s.CreateIngress(namespace, ingress)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth Ingress: %s\n", createdIngress.GetName())
	l2geth.PublicURL = publicURL
	_, err = t.db.UpdateComponent(l2geth)
	if err != nil {
		return
	}

	return
}
