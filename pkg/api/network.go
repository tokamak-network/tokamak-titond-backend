package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

func (t *TitondAPI) CreateNetwork(data *model.Network) (*model.Network, error) {
	result, err := t.db.CreateNetwork(data)
	if err == nil {
		go t.CreateNetworkInBackground(result)
	}
	return result, err
}

func (t *TitondAPI) CreateNetworkInBackground(network *model.Network) {
	namespace := t.ctx.String("titond.namespace")
	deployerName := MakeDeployerName(network.ID)
	t.k8s.CreateDeployer(namespace, deployerName)
	_ = t.k8s.WaitingDeploymentCreated(namespace, deployerName)
	podList, err := t.k8s.GetPodsOfDeployment(namespace, deployerName)
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
	addressData, dumpData := t.k8s.GetDeployerResult(namespace, &podList.Items[0])
	addressFileName := fmt.Sprintf("address-%d.json", network.ID)
	dumpFileName := fmt.Sprintf("state-dump-%d.json", network.ID)
	fmt.Println("Upload contract address file and genesis file")
	addressUrl, uploadAddressErr := t.UploadAddressFile(addressFileName, addressData)
	dumpUrl, uploadDumpErr := t.UploadDumpFile(dumpFileName, dumpData)

	err = t.UpdateDBWithValue(network, addressUrl, dumpUrl, uploadAddressErr, uploadDumpErr)
	if err == nil {
		fmt.Println("Clean k8s job", t.CleanK8sJob(network))
	}
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
	namespace := t.ctx.String("titond.namespace")
	deployerName := MakeDeployerName(network.ID)

	return t.k8s.DeleteDeployer(namespace, deployerName)
}

func MakeDeployerName(id uint) string {
	return fmt.Sprintf("deployer-%d", id)
}
