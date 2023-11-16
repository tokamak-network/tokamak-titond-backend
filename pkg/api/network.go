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
	t.k8s.CreateDeployer(namespace, "test2")
	_ = t.k8s.WaitingDeploymentCreated(namespace, "test2")
	podList, err := t.k8s.GetPodsOfDeployment(namespace, "test2")
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
	data := t.k8s.GetDeployerResult(namespace, &podList.Items[0])
	fmt.Println("data = ", data)
}
