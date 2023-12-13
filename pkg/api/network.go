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
	namespace := generateNamespace(network.ID)
	if _, err := t.k8s.CreateNamespace(namespace); err != nil {
		fmt.Printf("failed create namespace: %s\n", namespace)
		return "", ""
	}

	deployerName := MakeDeployerName(network.ID)
	_, err := t.createDeployer(namespace, deployerName)
	if err != nil {
		fmt.Println("failed create deployer...")
		return "", ""
	}
	podList, err := t.k8s.GetPodsOfDeployment(namespace, deployerName)
	if err != nil {
		fmt.Println("failed get pod of deployer...")
		return "", ""
	}
	if len(podList.Items) == 0 {
		fmt.Println("pod not exist...")
		return "", ""
	}

	if err := t.createAccounts(namespace); err != nil {
		fmt.Printf("failed create account in %s\n", namespace)
		return "", ""
	}

	fmt.Println("Getting file from pods...")
	addressData, addressErr := t.k8s.GetFileFromPod(namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/addresses.json")
	dumpData, dumpErr := t.k8s.GetFileFromPod(namespace, &podList.Items[0], "/opt/optimism/packages/tokamak/contracts/genesis/state-dump.latest.json")

	fmt.Printf("Get file contract address :\n%s\n", addressData)
	fmt.Printf("Get file dump data: \n%s\n", dumpData)
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
		fmt.Println("Clean k8s job")
		if err := t.k8s.DeleteDeployment(namespace, deployerName); err != nil {
			fmt.Println("Clean k8s job error: ", err)
		}
	}

	return addressUrl, dumpUrl
}

func (t *TitondAPI) createDeployer(namespace string, name string) (*appsv1.Deployment, error) {
	fmt.Println("Create deployer: ", name)
	mPath := t.k8s.GetManifestPath()

	obj := kubernetes.GetObject(mPath, "deployer", "configmap")
	cm, ok := kubernetes.ConvertToConfigMap(obj)
	if !ok {
		return nil, errors.New("createDeployer error: convertToConfigmap")
	}

	deployerConfig := map[string]string{
		"CONTRACTS_RPC_URL":        t.config.ContractsRpcUrl,
		"CONTRACTS_TARGET_NETWORK": t.config.ContractsTargetNetwork,
		"CONTRACTS_DEPLOYER_KEY":   t.config.ContractsDeployerKey,
	}

	_, err := t.k8s.CreateConfigMapWithConfig(namespace, cm, deployerConfig)
	if err != nil {
		return nil, errors.New("createDeployer error: " + err.Error())
	}

	obj = kubernetes.GetObject(mPath, "deployer", "deployment")
	deployment, success := kubernetes.ConvertToDeployment(obj)
	if !success {
		return nil, errors.New("createDeployer error: Failed to convert to deployment object")
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

func (t *TitondAPI) createAccounts(namespace string) error {
	stringData := map[string]string{
		"BATCH_SUBMITTER_SEQUENCER_PRIVATE_KEY": t.config.SequencerKey,
		"BATCH_SUBMITTER_PROPOSER_PRIVATE_KEY":  t.config.ProposerKey,
		"MESSAGE_RELAYER__L1_WALLET":            t.config.RelayerWallet,
		"BLOCK_SIGNER_KEY":                      t.config.SignerKey,
	}
	fmt.Printf("Create account in %s namespace\n", namespace)

	_, err := t.k8s.CreateSecret(namespace, "titan-secret", stringData)
	return err
}
