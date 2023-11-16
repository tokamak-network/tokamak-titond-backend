package api

import (
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
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
	t.k8s.CreateNamespaceForApp(t.ctx.String("titond.namespace"))
	t.k8s.CreateConfigMapForApp(t.ctx.String("titond.namespace"),
		t.ctx.String("titond.contracts.rpc.url"),
		t.ctx.String("titond.contracts.target.network"),
		t.ctx.String("titond.contracts.deployer.key"),
	)
}
