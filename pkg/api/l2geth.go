package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

type l2gethConfig struct {
	data map[string]string
}

func (t *TitondAPI) CreateL2Geth(data *model.Component /*TODO : get params for config, namespace*/) (*model.Component, error) {
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

	l2GethChan := make(chan error)

	go t.createL2Geth(namespace, config, l2GethChan)

	if err := <-l2GethChan; err != nil {
		return nil, err
	}
	close(l2GethChan)

	return data, nil
}

func (t *TitondAPI) createL2Geth(namespace string, config *l2gethConfig, res chan error) {
	obj := kubernetes.GetObject("l2geth", "configMap")
	cm, err := kubernetes.ConvertToConfigMap(obj)
	if err != nil {
		res <- err
		return
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, config.data)
	if err != nil {
		res <- err
		return
	}
	fmt.Printf("Created L1Geth ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject("l2geth", "pvc")
	pvc, err := kubernetes.ConvertToPersistentVolumeClaim(obj)
	if err != nil {
		res <- err
		return
	}

	createdPVC, err := t.k8s.CreatePersistentVolumeClaim(namespace, pvc)
	if err != nil {
		res <- err
		return
	}
	fmt.Printf("Created L1Geth PersistentVolumeClaim: %s\n", createdPVC.GetName())

	obj = kubernetes.GetObject("l2geth", "service")
	svc, err := kubernetes.ConvertToService(obj)
	if err != nil {
		res <- err
		return
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		res <- err
		return
	}
	fmt.Printf("Created L1Geth Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject("l2geth", "statefulset")
	sfs, err := kubernetes.ConvertToStatefulSet(obj)
	if err != nil {
		res <- err
		return
	}

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		res <- err
		return
	}
	fmt.Printf("Created L1Geth StatefulSet: %s\n", createdSFS.GetName())

	obj = kubernetes.GetObject("l2geth", "ingress")
	ingress, err := kubernetes.ConvertToIngress(obj)
	if err != nil {
		res <- err
		return
	}

	createdIngress, err := t.k8s.CreateIngress(namespace, ingress)
	if err != nil {
		res <- err
		return
	}
	fmt.Printf("Created L1Geth Ingress: %s\n", createdIngress.GetName())
	res <- nil
}
