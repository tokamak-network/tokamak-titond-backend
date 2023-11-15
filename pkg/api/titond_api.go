package api

import (
	"fmt"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/model"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
	"github.com/urfave/cli/v2"
)

type TitondAPI struct {
	k8s         *kubernetes.Kubernetes
	db          db.Client
	fileManager services.IFIleManager
	ctx         *cli.Context
}

func NewTitondAPI(k8s *kubernetes.Kubernetes, db db.Client, fileManager services.IFIleManager, ctx *cli.Context) *TitondAPI {
	return &TitondAPI{
		k8s,
		db,
		fileManager,
		ctx,
	}
}

func (t *TitondAPI) Initialize() {
	t.k8s.CreateConfigMapForApp(t.ctx.String("titond.namespace"),
		t.ctx.String("titond.contracts.rpc.url"),
		t.ctx.String("titond.contracts.deployer.key"),
		t.ctx.String("titond.contracts.target.network"),
	)
}

func (t *TitondAPI) CreateNetwork(data *model.Network) *model.Network {
	// t.fileManager.UploadContent("File_name_9", " New Content 9 ")
	result, _ := t.db.CreateNetwork(data)
	status, err := t.k8s.GetPodStatus(t.ctx.String("titond.namespace"), "l2geth-0")
	// status, err := t.k8s.GetPodStatus("default", "l2geth-0")

	fmt.Println("Create Network", status, " Err", err)

	return result
}
