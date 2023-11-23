package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

type l2gethConfig struct {
	data map[string]string
}

func (t *TitondAPI) CreateL2Geth(data *model.Component /*TODO : get params for config, namespace*/) *model.Component {
	// TODO : deal with DB
	// t.db.CreateComponent()

	/*
		This is currently hardcoding, but
		will be replaced by taking a value as a parameter and generating it
	*/
	namespace := "default"
	config := &l2gethConfig{
		data: map[string]string{},
	}

	config.data["ETH1_CONFIRMATION_DEPTH"] = "1"
	config.data["GASPRICE"] = "100"

	go t.createL2Geth(namespace, config)

	return data
}

func (t *TitondAPI) createL2Geth(namespace string, config *l2gethConfig) {
	obj := kubernetes.GetObject("l2geth", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToConfigmap")
		return
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, config.data)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject("l2geth", "pvc")
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

	obj = kubernetes.GetObject("l2geth", "service")
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

	obj = kubernetes.GetObject("l2geth", "statefulset")
	sfs, ok := kubernetes.ConvertToStatefulSet(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToStatefulSet")
		return
	}

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth StatefulSet: %s\n", createdSFS.GetName())

	obj = kubernetes.GetObject("l2geth", "ingress")
	ingress, ok := kubernetes.ConvertToIngress(obj)
	if !ok {
		fmt.Printf("createL2Geth error: convertToIngress")
		return
	}

	createdIngress, err := t.k8s.CreateIngress(namespace, ingress)
	if err != nil {
		fmt.Printf("createL2Geth error: %s\n", err)
		return
	}
	fmt.Printf("Created L2Geth Ingress: %s\n", createdIngress.GetName())
}
