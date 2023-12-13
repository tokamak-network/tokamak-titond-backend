package api

import (
	"errors"
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateExplorer(explorer *model.Component) (*model.Component, error) {
	network, err := t.db.ReadNetwork(explorer.NetworkID)
	if err != nil {
		return nil, err
	}

	l2geth, err := t.db.ReadComponentByType("l2geth", network.ID)
	if err != nil {
		return nil, err
	}
	if err := checkDependency(l2geth.Status); err != nil {
		return nil, err
	}

	result, err := t.db.CreateComponent(explorer)
	if err != nil {
		return nil, err
	}
	go t.createExplorer(result, network.StateDumpURL)

	return result, nil
}

func (t *TitondAPI) createExplorerDB(namespace string, explorer *model.Component) error {
	volumePath := generateVolumePath("explorer", explorer.NetworkID, explorer.ID)
	volumeLabel := generateLabel("explorer-pv", explorer.NetworkID, explorer.ID)

	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "volume", "pv")
	pv, ok := kubernetes.ConvertToPersistentVolume(obj)
	if !ok {
		return errors.New("createExplorerDB error: convertToPersistentVolume")
	}

	pv.SetName(volumeLabel)
	label := map[string]string{
		"app": volumeLabel,
	}

	createdPV, err := t.k8s.CreatePersistentVolume(label, "10Gi", pv)
	if err != nil {
		return err
	}
	fmt.Printf("Created Explorer DB PersistentVolume: %s\n", createdPV.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/postgresql", "pvc")
	pvc, ok := kubernetes.ConvertToPersistentVolumeClaim(obj)
	if !ok {
		return errors.New("createExplorerDB error: convertToPersistentVolumeClaim")
	}

	createdPVC, err := t.k8s.CreatePersistentVolumeClaim(namespace, label, "10Gi", pvc)
	if err != nil {
		return err
	}
	fmt.Printf("Created Explorer DB PersistentVolumeClaim: %s\n", createdPVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/postgresql", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		return errors.New("createExplorerDB error: convertToConfigmap")
	}

	createdConfigMap, err := t.k8s.CreateConfigMap(namespace, cm)
	if err != nil {
		return err
	}
	fmt.Printf("Created Explorer DB ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/postgresql", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		return errors.New("createExplorerDB error: convertToService")
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		return err
	}
	fmt.Printf("Created Explorer DB Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/postgresql", "statefulset")
	sfs, ok := kubernetes.ConvertToStatefulSet(obj)
	if !ok {
		return errors.New("createExplorerDB error: convertToStatefulSet")
	}

	sfs.Spec.Template.Spec.Containers[0].VolumeMounts[0].SubPath = volumePath

	createdSFS, err := t.k8s.CreateStatefulSet(namespace, sfs)
	if err != nil {
		return err
	}

	err = t.k8s.WatingStatefulsetCreated(namespace, createdSFS.Name)
	if err != nil {
		return err
	}
	fmt.Printf("Created Explorer DB StatefulSet: %s\n", createdSFS.GetName())

	return nil
}

func (t *TitondAPI) createSigProvider(namespace string) error {
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "explorer/sig-provider", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		return errors.New("createSigProvider error: convertToService")
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		return err
	}
	fmt.Printf("Created Sig-provider Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/sig-provider", "deployment")
	deployment, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		return errors.New("createSigProvider error: convertToDeployment")
	}

	createdDeplyment, err := t.k8s.CreateDeployment(namespace, deployment)
	if err != nil {
		return err
	}

	err = t.k8s.WaitingDeploymentCreated(namespace, createdDeplyment.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Created Sig-provide Deployment: %s\n", createdDeplyment.GetName())
	return nil
}
func (t *TitondAPI) createContractVerifier(namespace string) error {
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "explorer/smart-contract-verifier", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		return errors.New("createContractVerifier error: convertToConfigmap")
	}

	createdConfigMap, err := t.k8s.CreateConfigMap(namespace, cm)
	if err != nil {
		return err
	}
	fmt.Printf("Created Contract Verifier ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/smart-contract-verifier", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		return errors.New("createContractVerifier error: convertToService")
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		return err
	}
	fmt.Printf("Created Contract Verifier Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/smart-contract-verifier", "deployment")
	deployment, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		return errors.New("createContractVerifier error: convertToDeployment")
	}

	createdDeplyment, err := t.k8s.CreateDeployment(namespace, deployment)
	if err != nil {
		return err
	}

	err = t.k8s.WaitingDeploymentCreated(namespace, createdDeplyment.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Created Contract Verifier Deployment: %s\n", createdDeplyment.GetName())

	return nil
}
func (t *TitondAPI) createVisualizer(namespace string) error {
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "explorer/visualizer", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		return errors.New("createVisualizer error: convertToConfigmap")
	}

	createdConfigMap, err := t.k8s.CreateConfigMap(namespace, cm)
	if err != nil {
		return err
	}
	fmt.Printf("Created Visualizer ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/visualizer", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		return errors.New("createVisualizer error: convertToService")
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		return err
	}
	fmt.Printf("Created Visualizer Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer/visualizer", "deployment")
	deployment, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		return errors.New("createVisualizer error: convertToDeployment")
	}

	createdDeplyment, err := t.k8s.CreateDeployment(namespace, deployment)
	if err != nil {
		return err
	}

	err = t.k8s.WaitingDeploymentCreated(namespace, createdDeplyment.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Created Visualizer Deployment: %s\n", createdDeplyment.GetName())

	return nil
}

func (t *TitondAPI) createExplorer(explorer *model.Component, stateDumpURL string) {
	namespace := generateNamespace(explorer.NetworkID)
	publicURL := generatePublcURL("explorer", explorer.NetworkID, explorer.ID)

	if err := t.createExplorerDB(namespace, explorer); err != nil {
		fmt.Printf("createExplorer error: %v\n", err)
		return
	}

	if err := t.createSigProvider(namespace); err != nil {
		fmt.Printf("createExplorer error: %v\n", err)
		return
	}

	if err := t.createContractVerifier(namespace); err != nil {
		fmt.Printf("createExplorer error: %v\n", err)
		return
	}

	if err := t.createVisualizer(namespace); err != nil {
		fmt.Printf("createExplorer error: %v\n", err)
		return
	}

	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "explorer", "configMap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		fmt.Printf("createExplorer error: convertToConfigmap")
		return
	}

	explorerConfig := map[string]string{
		"CHAIN_SPEC_PATH": stateDumpURL,
	}

	createdConfigMap, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, explorerConfig)
	if err != nil {
		fmt.Printf("createExplorer error: %s\n", err)
		return
	}
	fmt.Printf("Created Explorer ConfigMap: %s\n", createdConfigMap.GetName())

	obj = kubernetes.GetObject(mPath, "explorer", "service")
	svc, ok := kubernetes.ConvertToService(obj)
	if !ok {
		fmt.Printf("createExplorer error: convertToService")
		return
	}

	createdSVC, err := t.k8s.CreateService(namespace, svc)
	if err != nil {
		fmt.Printf("createExplorer error: %s\n", err)
		return
	}
	fmt.Printf("Created Explorer Service: %s\n", createdSVC.GetName())

	obj = kubernetes.GetObject(mPath, "explorer", "deployment")
	deployment, ok := kubernetes.ConvertToDeployment(obj)
	if !ok {
		fmt.Printf("createExplorer error: convertToDeployment")
		return
	}

	createdDeplyment, err := t.k8s.CreateDeployment(namespace, deployment)
	if err != nil {
		fmt.Printf("createExplorer error: %s\n", err)
		return
	}
	fmt.Printf("Created Explorer StatefulSet: %s\n", createdDeplyment.GetName())

	err = t.k8s.WaitingDeploymentCreated(createdDeplyment.Namespace, createdDeplyment.Name)
	if err != nil {
		/*TODO : rollback ? */
		return
	}
	explorer.Status = true

	obj = kubernetes.GetObject(mPath, "explorer", "ingress")
	ingress, ok := kubernetes.ConvertToIngress(obj)
	if !ok {
		fmt.Printf("createExplorer error: convertToIngress")
		return
	}

	ingress.Spec.Rules[0].Host = publicURL
	ingress.Spec.TLS[0].Hosts[0] = publicURL

	createdIngress, err := t.k8s.CreateIngress(namespace, ingress)
	if err != nil {
		fmt.Printf("createExplorer error: %s\n", err)
		return
	}
	fmt.Printf("Created Explorer Ingress: %s\n", createdIngress.GetName())
	explorer.PublicURL = publicURL
	_, err = t.db.UpdateComponent(explorer)
	if err != nil {
		return
	}

	return
}
