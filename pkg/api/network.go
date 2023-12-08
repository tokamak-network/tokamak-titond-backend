package api

import (
	"errors"
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
)

var PAGE_SIZE int = 10

func (t *TitondAPI) CreateNetwork(data *model.Network) (*model.Network, error) {
	result, err := t.db.CreateNetwork(data)
	if err == nil {
		go t.createNetwork(result)
	}
	return result, err
}

func (t *TitondAPI) GetNetworksByPage(page int) ([]model.Network, error) {
	networks, err := t.db.ReadNetworkByRange((page-1)*PAGE_SIZE, PAGE_SIZE)
	if err == nil {
		if len(networks) == 0 {
			return nil, types.ErrResourceNotFound
		}
	}
	return networks, err
}

func (t *TitondAPI) GetNetworkByID(networkID uint) (*model.Network, error) {
	return t.db.ReadNetwork(networkID)
}

func (t *TitondAPI) DeleteNetwork(id uint) error {
	result, err := t.db.DeleteNetwork(id)
	if err == nil {
		if result == 0 {
			return types.ErrResourceNotFound
		}
	}
	return err
}

func (t *TitondAPI) createNetwork(network *model.Network) (string, string) {
	deployerName := MakeDeployerName(network.ID)
	_, err := t.createDeployer(t.config.Namespace, deployerName)
	if err != nil {
		return "", ""
	}
	podList, err := t.k8s.GetPodsOfDeployment(t.config.Namespace, deployerName)
	if err != nil {
		return "", ""
	}
	if len(podList.Items) == 0 {
		return "", ""
	}

	namespace := generateNamespace(network.ID)
	if _, err := t.k8s.CreateNamespace(namespace); err != nil {
		return "", ""
	}
	if err := t.createAccounts(namespace); err != nil {
		return "", ""
	}

	addressData, addressErr := t.k8s.GetFileFromPod(t.config.Namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/addresses.json")
	dumpData, dumpErr := t.k8s.GetFileFromPod(t.config.Namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json")

	addressUrl := ""
	dumpUrl := ""
	uploadAddressErr := errors.New("")
	uploadDumpErr := errors.New("")

	if addressErr == nil {
		addressFileName := fmt.Sprintf("address-%d.json", network.ID)
		addressUrl, uploadAddressErr = t.uploadAddressFile(addressFileName, addressData)
	}
	if dumpErr == nil {
		dumpFileName := fmt.Sprintf("state-dump-%d.json", network.ID)
		dumpUrl, uploadDumpErr = t.uploadDumpFile(dumpFileName, dumpData)
	}
	err = t.updateDBWithValue(network, addressUrl, dumpUrl, uploadAddressErr, uploadDumpErr)

	// We clear the job when everything works correctly
	if err == nil && uploadAddressErr == nil && uploadDumpErr == nil {
		fmt.Println("Clean k8s job", t.cleanK8sJob(network))
	}

	return addressUrl, dumpUrl
}

func (t *TitondAPI) createDeployer(namespace string, name string) (*appsv1.Deployment, error) {
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

func (t *TitondAPI) uploadAddressFile(addressFileName string, addressData string) (string, error) {
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

func (t *TitondAPI) uploadDumpFile(dumpFileName string, dumpFileData string) (string, error) {
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

func (t *TitondAPI) updateDBWithValue(network *model.Network, addressFileUrl string, dumpFileUrl string, uploadAddressErr error, uploadDumpErr error) error {
	network.ContractAddressURL = addressFileUrl
	network.StateDumpURL = dumpFileUrl
	network.Status = (uploadAddressErr == nil) && (uploadDumpErr == nil)
	_, err := t.db.UpdateNetwork(network)
	return err
}

func (t *TitondAPI) cleanK8sJob(network *model.Network) error {
	deployerName := MakeDeployerName(network.ID)
	return t.k8s.DeleteDeployment(t.config.Namespace, deployerName)
}

func (t *TitondAPI) GetK8sJobStatus(network *model.Network) (*appsv1.Deployment, error) {
	deployerName := MakeDeployerName(network.ID)
	return t.k8s.GetDeployment(t.config.Namespace, deployerName)
}

func (t *TitondAPI) createAccounts(namespace string) error {
	sequencerKey, address := generateKey()
	fmt.Printf("created sequencer account: %s\n", address)

	proposerKey, address := generateKey()
	fmt.Printf("created proposer account: %s\n", address)

	relayerKey, address := generateKey()
	fmt.Printf("created relayer account: %s\n", address)

	signerKey, address := generateKey()
	fmt.Printf("created block signer account: %s\n", address)

	stringData := map[string]string{
		"BATCH_SUBMITTER_SEQUENCER_PRIVATE_KEY": sequencerKey,
		"BATCH_SUBMITTER_PROPOSER_PRIVATE_KEY":  proposerKey,
		"MESSAGE_RELAYER__L1_WALLET":            relayerKey,
		"BLOCK_SIGNER_KEY":                      signerKey,
	}

	_, err := t.k8s.CreateSecret(namespace, "titan-secret", stringData)
	return err
}
