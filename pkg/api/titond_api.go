package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
)

type TitondAPI struct {
	k8s *kubernetes.Kubernetes
	db  db.Client
}

func NewTitondAPI(k8s *kubernetes.Kubernetes, db db.Client) *TitondAPI {
	return &TitondAPI{
		k8s,
		db,
	}
}

func (t *TitondAPI) CreateNetwork(data *model.Network) *model.Network {
	result, _ := t.db.CreateNetwork(data)
	status, _ := t.k8s.GetPodStatus("default", "l2geth-0")
	fmt.Println(status)

	return result
}
