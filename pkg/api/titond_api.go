package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
)

type TitondAPI struct {
	k8s         *kubernetes.Kubernetes
	db          db.Client
	fileManager services.IFIleManager
}

func NewTitondAPI(k8s *kubernetes.Kubernetes, db db.Client, fileManager services.IFIleManager) *TitondAPI {
	return &TitondAPI{
		k8s,
		db,
		fileManager,
	}
}

func (t *TitondAPI) CreateNetwork(data *model.Network) *model.Network {
	// t.fileManager.UploadContent("File_name_2", " New Content 2 ")
	result, _ := t.db.CreateNetwork(data)
	status, _ := t.k8s.GetPodStatus("default", "l2geth-0")
	fmt.Println(status)

	return result
}
