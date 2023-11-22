package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
)

func (t *TitondAPI) CreateNetwork(data *model.Network) (*model.Network, error) {
	result, err := t.db.CreateNetwork(data)
	if err == nil {
		go t.createNetwork(result)
	}
	return result, err
}

func (t *TitondAPI) createNetwork(network *model.Network) {
	deployerName := MakeDeployerName(network.ID)
	_, err := t.CreateDeployer(t.config.Namespace, deployerName)
	if err != nil {
		fmt.Println("Failed when creating deployer:", err)
		return
	}
	podList, err := t.k8s.GetPodsOfDeployment(t.config.Namespace, deployerName)
	if err != nil {
		return
	}
	if podList != nil {
		fmt.Println("Pod len", len(podList.Items))
		for _, pod := range podList.Items {
			fmt.Println("Pod name: ", pod.Name)
		}
	}
	if len(podList.Items) == 0 {
		fmt.Println("Back")
		return
	}
	// addressData, addressErr, dumpData, dumpErr := t.k8s.GetDeployerResult(t.config.Namespace, &podList.Items[0])
	addressData, addressErr := t.k8s.GetFileFromPod(t.config.Namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/addresses.json")
	dumpData, dumpErr := t.k8s.GetFileFromPod(t.config.Namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json")

	addressUrl := ""
	dumpUrl := ""
	var uploadDumpErr, uploadAddressErr error
	if addressErr == nil {
		addressFileName := fmt.Sprintf("address-%d.json", network.ID)
		addressUrl, uploadAddressErr = t.UploadAddressFile(addressFileName, addressData)
	}
	if dumpErr == nil {
		dumpFileName := fmt.Sprintf("state-dump-%d.json", network.ID)
		dumpUrl, uploadDumpErr = t.UploadDumpFile(dumpFileName, dumpData)
	}
	err = t.UpdateDBWithValue(network, addressUrl, dumpUrl, uploadAddressErr, uploadDumpErr)
	if err == nil {
		fmt.Println("Clean k8s job", t.CleanK8sJob(network))
	}
}

func (t *TitondAPI) CreateDeployer(namespace string, name string) (*appsv1.Deployment, error) {
	fmt.Println("Create deployer: ", name)
	object, _ := kubernetes.BuildObjectFromYamlFile("./deployments/deployer/deployment.yaml")
	deployment, success := kubernetes.ConvertToDeployment(object)
	if !success {
		panic("Failed to convert to deployment object")
	}
	var deployerCreationErr error
	_, deployerCreationErr = t.k8s.CreateDeploymentWithName(namespace, deployment, name)
	if deployerCreationErr == nil {
		deployerCreationErr = t.k8s.WaitingDeploymentCreated(namespace, name)
	}

	return deployment, deployerCreationErr
}

func (t *TitondAPI) UploadAddressFile(addressFileName string, addressData string) (string, error) {
	var addressUrl string
	var err error
	for i := 0; i < 5; i++ {
		addressUrl, err = t.fileManager.UploadContent(addressFileName, addressData)
		if err == nil {
			return addressUrl, err
		}
	}
	return "", err
}

func (t *TitondAPI) UploadDumpFile(dumpFileName string, dumpFileData string) (string, error) {
	var dumpFileUrl string
	var err error
	for i := 0; i < 5; i++ {
		dumpFileUrl, err = t.fileManager.UploadContent(dumpFileName, dumpFileData)
		if err == nil {
			return dumpFileUrl, err
		}
	}
	return "", err
}

func (t *TitondAPI) UpdateDBWithValue(network *model.Network, addressFileUrl string, dumpFileUrl string, uploadAddressErr error, uploadDumpErr error) error {
	network.ContractAddressURL = addressFileUrl
	network.StateDumpURL = dumpFileUrl
	network.Status = (uploadAddressErr == nil) && (uploadDumpErr == nil)
	_, err := t.db.UpdateNetwork(network)
	return err
}

func (t *TitondAPI) CleanK8sJob(network *model.Network) error {
	deployerName := MakeDeployerName(network.ID)

	return t.k8s.DeleteDeployment(t.config.Namespace, deployerName)
}

func MakeDeployerName(id uint) string {
	return fmt.Sprintf("deployer-%d", id)
}
